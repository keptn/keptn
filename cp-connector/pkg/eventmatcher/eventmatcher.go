package eventmatcher

import (
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/sliceutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

// EventMatcher is used to check whether an event contains is containing information
// about a specif event, stage or service
type EventMatcher struct {
	Project string
	Stage   string
	Service string
}

// New creates a new EventMatcher that is configured
// with information about project, stage and service filter contained in an event subscription
func New(subscription models.EventSubscription) *EventMatcher {
	return &EventMatcher{
		Project: strings.Join(subscription.Filter.Projects, ","),
		Stage:   strings.Join(subscription.Filter.Stages, ","),
		Service: strings.Join(subscription.Filter.Services, ","),
	}
}

// Matches checks whether a Keptn event matches the information of the currently configured
// EventMatcher
func (ef EventMatcher) Matches(e models.KeptnContextExtendedCE) bool {
	generalEventData := &v0_2_0.EventData{}
	if err := e.DataAs(generalEventData); err != nil {
		return false
	}

	if ef.Project != "" && !sliceutils.ContainsStr(strings.Split(ef.Project, ","), generalEventData.Project) ||
		ef.Stage != "" && !sliceutils.ContainsStr(strings.Split(ef.Stage, ","), generalEventData.Stage) ||
		ef.Service != "" && !sliceutils.ContainsStr(strings.Split(ef.Service, ","), generalEventData.Service) {
		return false
	}
	return true
}
