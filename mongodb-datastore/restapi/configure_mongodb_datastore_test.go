package restapi

import (
	"context"
	"github.com/keptn/keptn/mongodb-datastore/handlers"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations"
	"github.com/nats-io/nats-server/v2/server"
	natstest "github.com/nats-io/nats-server/v2/test"
	logger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

// Test_startControlPlane verifies whether the api pre server shutdown is initialized
// and that the control plane terminates after calling it
func Test_startControlPlaneSuccess(t *testing.T) {
	natsServer, shutdown := func() (*server.Server, func()) {
		svr := natstest.RunRandClientPortServer()
		return svr, func() { svr.Shutdown() }
	}()
	defer func() {
		shutdown()
	}()
	err := os.Setenv("NATS_URL", natsServer.ClientURL())
	require.NoError(t, err)

	api := &operations.MongodbDatastoreAPI{}
	eventRequestHandler := handlers.NewEventRequestHandler(nil)
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		err := startControlPlane(ctx, api, eventRequestHandler, logger.New())
		require.Nil(t, err)
		t.Log("control plane terminated")
	}()
	// test propagate shutdown
	require.Eventually(t, func() bool {
		return getPreShutDown(api) != nil
	}, 10*time.Second, 1*time.Second)
	cancel()

}

func Test_startControlPlaneFailNoNATS(t *testing.T) {
	eventRequestHandler := handlers.NewEventRequestHandler(nil)

	log := &logger.Logger{
		Out:       os.Stderr,
		Formatter: new(logger.TextFormatter),
		Hooks:     make(logger.LevelHooks),
		Level:     logger.DebugLevel,
	}
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		time.Sleep(10 * time.Second)
		cancel()
	}()
	err := startControlPlane(ctx, &operations.MongodbDatastoreAPI{}, eventRequestHandler, log)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Could not handle subscription update")

}
