package main

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	fakek8s "k8s.io/client-go/kubernetes/fake"
	fakeclient "k8s.io/client-go/testing"
	"os"
	"testing"
	"time"
)

func Test_getDurationFromEnvVar(t *testing.T) {
	type args struct {
		envVarValue string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "get default value",
			args: args{
				envVarValue: "",
			},
			want: 432000 * time.Second,
		},
		{
			name: "get configured value",
			args: args{
				envVarValue: "10s",
			},
			want: 10 * time.Second,
		},
		{
			name: "get configured value",
			args: args{
				envVarValue: "2m",
			},
			want: 120 * time.Second,
		},
		{
			name: "get configured value",
			args: args{
				envVarValue: "1h30m",
			},
			want: 5400 * time.Second,
		},
		{
			name: "get default value because of invalid config",
			args: args{
				envVarValue: "invalid",
			},
			want: 432000 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LOG_TTL", tt.args.envVarValue)
			if got := getDurationFromEnvVar("LOG_TTL", envVarLogsTTLDefault); got != tt.want {
				t.Errorf("getLogTTLDurationInSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_LeaderElection(t *testing.T) {
	var (
		onNewLeader = make(chan struct{})
		onRelease   = make(chan struct{})
		lockObj     runtime.Object
	)

	c := &fakek8s.Clientset{}

	shipyard := &fake.IShipyardControllerMock{
		StartDispatchersFunc: func(ctx context.Context) {
			close(onNewLeader)
		},
		StopDispatchersFunc: func() {
			onNewLeader = make(chan struct{})
			close(onRelease)
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	// create lock
	c.AddReactor("create", "leases", func(action fakeclient.Action) (handled bool, ret runtime.Object, err error) {
		lockObj = action.(fakeclient.CreateAction).GetObject()
		return true, lockObj, nil
	})

	//fail with no lock
	c.AddReactor("get", "leases", func(action fakeclient.Action) (handled bool, ret runtime.Object, err error) {
		if lockObj != nil {
			return true, lockObj, nil
		}
		return true, nil, errors.NewNotFound(action.(fakeclient.GetAction).GetResource().GroupResource(), action.(fakeclient.GetAction).GetName())
	})

	c.AddReactor("update", "leases", func(action fakeclient.Action) (handled bool, ret runtime.Object, err error) {
		// Second update (first renew) should return our canceled error
		// FakeClient doesn't do anything with the context so we're doing this ourselves

		lockObj = action.(fakeclient.UpdateAction).GetObject()
		return true, lockObj, nil

	})

	c.AddReactor("*", "*", func(action fakeclient.Action) (bool, runtime.Object, error) {
		t.Errorf("unreachable action. testclient called too many times: %+v", action)
		return true, nil, fmt.Errorf("unreachable action")
	})

	newReplica := func() { LeaderElection(c.CoordinationV1(), ctx, shipyard.StartDispatchers, shipyard.StopDispatchers) }
	go newReplica()

	// Wait for first replica to become the leader
	select {
	case <-onNewLeader:
	case <-time.After(10 * time.Second):
		t.Fatal("failed to become the leader")
	}

	// leader already there this part should fail but not panic
	go newReplica()

	time.After(25 * time.Second)
	// stopping the leader
	cancel()

	select {
	case <-onRelease:
	case <-time.After(10 * time.Second):
		t.Fatal("the lock was not released")
	}
}
