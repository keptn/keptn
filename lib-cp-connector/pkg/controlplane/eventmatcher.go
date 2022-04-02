package controlplane

import (
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/sliceutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"strings"
)

type EventMatcher struct {
	Project string
	Stage   string
	Service string
}

func NewEventMatcherFromSubscription(subscription models.EventSubscription) *EventMatcher {
	return &EventMatcher{
		Project: strings.Join(subscription.Filter.Projects, ","),
		Stage:   strings.Join(subscription.Filter.Stages, ","),
		Service: strings.Join(subscription.Filter.Services, ","),
	}
}

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
