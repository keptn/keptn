package leaderelection

import (
	"context"
	"github.com/google/uuid"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/sirupsen/logrus"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/coordination/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"time"
)

func LeaderElection(client v1.CoordinationV1Interface, ctx context.Context, start func(ctx context.Context, mode common.SDMode), stop func()) {
	myID := uuid.New().String()
	// we use the Lease lock type since edits to Leases are less common
	// and fewer objects in the cluster watch "all Leases".
	lock := &resourcelock.LeaseLock{
		LeaseMeta: v12.ObjectMeta{
			Name:      "shipyard-controller-dispatcher",
			Namespace: common.GetKeptnNamespace(),
		},
		Client: client,
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: myID,
		},
	}

	// start the leader election code loop
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock: lock,
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				// we're notified when we start - this is where you would
				// usually put your code
				start(ctx, common.SDModeRW)
			},
			OnStoppedLeading: func() {
				// we can do cleanup here
				logrus.Infof("leader lost: %s", myID)
				stop()
			},
			OnNewLeader: func(identity string) {
				// we're notified when a new leader is elected
				if identity == myID {
					// I just got the lock
					return
				}
				logrus.Infof("new leader elected: %s", identity)
				stop()
			},
		},
	})
}
