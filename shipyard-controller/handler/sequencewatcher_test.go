package handler_test

import (
	"context"
	"github.com/benbjohnson/clock"
	"github.com/keptn/keptn/shipyard-controller/common"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"testing"
	"time"
)

func TestSequenceWatcher(t *testing.T) {
	watcher := handler.NewSequenceWatcher(
		&fake.IShipyardControllerMock{
			CancelSequenceFunc: func(cancelRequest handler.SequenceCancellation) error {
				return nil
			},
		},
		&db_mock.EventRepoMock{
			DeleteEventFunc: func(project string, eventID string, status common.EventStatus) error {
				return nil
			},
			GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
				return []models.Event{}, nil
			},
		},
		&db_mock.ProjectRepoMock{
			GetProjectsFunc: func() ([]*models.ExpandedProject, error) {
				return []*models.ExpandedProject{}, nil
			},
		},
		10*time.Minute,
		1*time.Minute,
		clock.NewMock(),
	)
	ctx, cancel := context.WithCancel(context.Background())

	watcher.Run(ctx)

	cancel()
}
