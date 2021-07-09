package events

import (
	"context"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/common/sliceutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

// NATSEventReceiver receives events directly from the NATS broker and sends the cloud
type NATSEventReceiver struct {
	env       config.EnvConfig
	ceClient  cloudevents.Client
	closeChan chan bool
}

func NewNATSEventReceiver(ceClient cloudevents.Client, env config.EnvConfig) *NATSEventReceiver {
	return &NATSEventReceiver{
		env:       env,
		ceClient:  ceClient,
		closeChan: make(chan bool),
	}
}

func (n *NATSEventReceiver) Start(ctx *ExecutionContext) {
	if n.env.PubSubRecipient == "" {
		logger.Warn("No pubsub recipient defined")
		return
	}
	if n.env.PubSubTopic == "" {
		logger.Warn("No pubsub topic defined. No need to create NATS client connection.")
		return
	}
	uptimeTicker := time.NewTicker(10 * time.Second)

	natsURL := n.env.PubSubURL

	topics := strings.Split(n.env.PubSubTopic, ",")
	nch := NewNatsConnectionHandler(natsURL, topics)

	nch.MessageHandler = n.handleMessage

	err := nch.SubscribeToTopics()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func() {
		nch.RemoveAllSubscriptions()
		logger.Info("Disconnected from NATS")
	}()

	for {
		select {
		case <-uptimeTicker.C:
			_ = nch.SubscribeToTopics()
		case <-n.closeChan:
			return
		case <-ctx.Done():
			logger.Info("Terminating NATS event receiver")
			ctx.Wg.Done()
			return
		}
	}
}

func (n *NATSEventReceiver) handleMessage(m *nats.Msg) {
	go func() {
		logger.Infof("Received a message for topic [%s]\n", m.Subject)
		e, err := DecodeCloudEvent(m.Data)

		if e != nil && err == nil {
			err = n.sendEvent(*e)
			if err != nil {
				logger.Errorf("Could not send CloudEvent: %v", err)
			}
		}
	}()
}

func (n *NATSEventReceiver) sendEvent(event cloudevents.Event) error {
	if !n.matchesFilter(event) {
		// Do not send cloud event if it does not match the filter
		return nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	ctx = cloudevents.ContextWithTarget(ctx, config.GetPubSubRecipientURL(n.env))
	ctx = cloudevents.WithEncodingStructured(ctx)
	defer cancel()

	if result := n.ceClient.Send(ctx, event); cloudevents.IsUndelivered(result) {
		fmt.Printf("failed to send: %s\n", result.Error())
		return errors.New(result.Error())
	}
	fmt.Printf("sent: %s\n", event.ID())
	return nil
}

func (n *NATSEventReceiver) matchesFilter(e cloudevents.Event) bool {
	keptnBase := &v0_2_0.EventData{}
	if err := e.DataAs(keptnBase); err != nil {
		return true
	}
	if n.env.ProjectFilter != "" && !sliceutils.ContainsStr(strings.Split(n.env.ProjectFilter, ","), keptnBase.Project) ||
		n.env.StageFilter != "" && !sliceutils.ContainsStr(strings.Split(n.env.StageFilter, ","), keptnBase.Stage) ||
		n.env.ServiceFilter != "" && !sliceutils.ContainsStr(strings.Split(n.env.ServiceFilter, ","), keptnBase.Service) {
		return false
	}
	return true
}
