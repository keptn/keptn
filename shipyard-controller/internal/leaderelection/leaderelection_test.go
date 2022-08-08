package leaderelection

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/internal/controller/fake"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	fakek8s "k8s.io/client-go/kubernetes/fake"
	fakeclient "k8s.io/client-go/testing"
	"testing"
	"time"
)

func Test_LeaderElection(t *testing.T) {
	// TODO need to move this to another place
	var (
		onNewLeader = make(chan struct{})
		onRelease   = make(chan struct{})
		lockObj     runtime.Object
	)
	c := &fakek8s.Clientset{}

	shipyard := &fake.IShipyardControllerMock{
		StartDispatchersFunc: func(ctx context.Context, mode common.SDMode) {
			time.After(10 * time.Second)
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

	newReplica := func() {
		LeaderElection(c.CoordinationV1(), ctx, shipyard.StartDispatchersFunc, shipyard.StopDispatchers)
	}
	go newReplica()

	// Wait for one replica to become the leader
	select {
	case <-onNewLeader:
		// stopping the leader

		go newReplica() // leader already there one of the two may fail but not panic
		cancel()
		select {
		case <-onRelease:
			//reset chan for next leader
			onRelease = make(chan struct{})
		case <-time.After(10 * time.Second):
			t.Fatal("failed to release lock")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("failed to become the leader")
	}
	cancel()

	require.Eventually(t, func() bool {
		return len(shipyard.StopDispatchersCalls()) > 0
	}, 5*time.Second, 100*time.Millisecond)
}
