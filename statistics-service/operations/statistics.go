package operations

import (
	"time"
)

// GetStatisticsParams godoc
type GetStatisticsParams struct {
	// From godoc
	From time.Time `form:"from" json:"from" time_format:"unix"`
	// To godoc
	To time.Time `form:"to" json:"to" time_format:"unix"`
}

// GetStatisticsResponse godoc
type GetStatisticsResponse struct {
	// From godoc
	From time.Time `json:"from" bson:"from"`
	// To godoc
	To time.Time `json:"to" bson:"to"`
	// Projects godoc
	Projects []GetStatisticsResponseProject `json:"projects" bson:"projects"`
}

// GetStatisticsResponseProject godoc
type GetStatisticsResponseProject struct {
	// Name godoc
	Name string `json:"name" bson:"name"`
	// Services godoc
	Services []GetStatisticsResponseService `json:"services" bson:"services"`
}

// GetStatisticsResponseService godoc
type GetStatisticsResponseService struct {
	// Name godoc
	Name string `json:"name" bson:"name"`
	// Events godoc
	Events []GetStatisticsResponseEvent `json:"events" bson:"events"`
	// KeptnServiceExecutions godoc
	KeptnServiceExecutions []GetStatisticsResponseKeptnService `json:"keptnServiceExecutions" bson:"keptnServiceExecutions"`
	// ExecutedSequencesPerType godoc
	ExecutedSequencesPerType []GetStatisticsResponseEvent `json:"executedSequencesPerType,omitempty" bson:"executedSequencesPerType"`
}

// GetStatisticsResponseEvent godoc+
type GetStatisticsResponseEvent struct {
	// Type godoc
	Type string `json:"type" bson:"type"`
	// Count
	Count int `json:"count" bson:"count"`
}

// GetStatisticsResponseKeptnService godoc
type GetStatisticsResponseKeptnService struct {
	// Name godoc
	Name string `json:"name" bson:"name"`
	// Executions godoc
	Executions []GetStatisticsResponseEvent `json:"executions" bson:"executions"`
}

// Statistics godoc
type Statistics struct {
	// From godoc
	From time.Time `json:"from" bson:"from"`
	// To godoc
	To time.Time `json:"to" bson:"to"`
	// Projects godoc
	Projects map[string]*Project `json:"projects" bson:"projects"`
}

// Project godoc
type Project struct {
	// Name godoc
	Name string `json:"name" bson:"name"`
	// Services godoc
	Services map[string]*Service `json:"services" bson:"services"`
}

// Service godoc
type Service struct {
	// Name godoc
	Name string `json:"name" bson:"name"`
	// ExecutedSequences godoc
	ExecutedSequences int `json:"executedSequences" bson:"executedSequences"`
	// ExecutedSequencesPerType godoc
	ExecutedSequencesPerType map[string]int `json:"executedSequencesPerType" bson:"executedSequencesPerType"`
	// Events godoc
	Events map[string]int `json:"events" bson:"events"`
	// KeptnServiceExecutions godoc
	KeptnServiceExecutions map[string]*KeptnService `json:"keptnServiceExecutions" bson:"keptnServiceExecutions"`
}

// KeptnService godoc
type KeptnService struct {
	// Name godoc
	Name string `json:"name" bson:"name"`
	// Executions godoc
	Executions map[string]int `json:"executions" bson:"executions"`
}

func (s *Statistics) ensureProjectAndServiceExist(projectName string, serviceName string) {
	s.ensureProjectExists(projectName)
	if s.Projects[projectName].Services == nil {
		s.Projects[projectName].Services = map[string]*Service{}
	}
	if s.Projects[projectName].Services[serviceName] == nil {
		s.Projects[projectName].Services[serviceName] = &Service{
			Name:                     serviceName,
			ExecutedSequences:        0,
			ExecutedSequencesPerType: map[string]int{},
			Events:                   map[string]int{},
			KeptnServiceExecutions:   map[string]*KeptnService{},
		}
	}
}

func (s *Statistics) ensureKeptnServiceExists(projectName, serviceName, keptnServiceName string) {
	s.ensureProjectAndServiceExist(projectName, serviceName)
	if s.Projects[projectName].Services[serviceName].KeptnServiceExecutions == nil {
		s.Projects[projectName].Services[serviceName].KeptnServiceExecutions = map[string]*KeptnService{}
	}
	if s.Projects[projectName].Services[serviceName].KeptnServiceExecutions[keptnServiceName] == nil {
		s.Projects[projectName].Services[serviceName].KeptnServiceExecutions[keptnServiceName] = &KeptnService{
			Name:       keptnServiceName,
			Executions: map[string]int{},
		}
	}
}

func (s *Statistics) ensureProjectExists(projectName string) {
	if s.Projects == nil {
		s.Projects = map[string]*Project{}
	}
	if s.Projects[projectName] == nil {
		s.Projects[projectName] = &Project{
			Name:     projectName,
			Services: map[string]*Service{},
		}
	}
}

// IncreaseEventTypeCount godoc
func (s *Statistics) IncreaseEventTypeCount(projectName, serviceName, eventType string, increment int) {
	s.ensureProjectAndServiceExist(projectName, serviceName)
	service := s.Projects[projectName].Services[serviceName]
	service.Events[eventType] = service.Events[eventType] + increment
}

// IncreaseExecutedSequencesCount godoc
func (s *Statistics) IncreaseExecutedSequencesCount(projectName, serviceName string, increment int) {
	s.ensureProjectAndServiceExist(projectName, serviceName)
	service := s.Projects[projectName].Services[serviceName]
	service.ExecutedSequences = service.ExecutedSequences + increment
}

// IncreaseKeptnServiceExecutionCount godoc
func (s *Statistics) IncreaseKeptnServiceExecutionCount(projectName, serviceName, keptnServiceName, eventType string, increment int) {
	s.ensureProjectAndServiceExist(projectName, serviceName)
	s.ensureKeptnServiceExists(projectName, serviceName, keptnServiceName)
	keptnService := s.Projects[projectName].Services[serviceName].KeptnServiceExecutions[keptnServiceName]
	keptnService.Executions[eventType] = keptnService.Executions[eventType] + increment
}

// IncreaseExecutedSequenceCountForType godoc
func (s *Statistics) IncreaseExecutedSequenceCountForType(projectName string, serviceName string, eventType string, increment int) {
	s.ensureProjectAndServiceExist(projectName, serviceName)
	service := s.Projects[projectName].Services[serviceName]
	service.ExecutedSequencesPerType[eventType] = service.ExecutedSequencesPerType[eventType] + increment
}

// MergeStatistics godoc
func MergeStatistics(target Statistics, statistics []Statistics) Statistics {
	for _, stats := range statistics {
		for projectName, project := range stats.Projects {
			target.ensureProjectExists(projectName)
			for serviceName, service := range project.Services {
				for eventType, count := range service.Events {
					target.IncreaseEventTypeCount(projectName, serviceName, eventType, count)
				}
				if service.ExecutedSequences > 0 {
					target.IncreaseExecutedSequencesCount(projectName, serviceName, service.ExecutedSequences)
				}
				for keptnServiceName, keptnService := range service.KeptnServiceExecutions {
					for eventType, count := range keptnService.Executions {
						target.IncreaseKeptnServiceExecutionCount(projectName, serviceName, keptnServiceName, eventType, count)
					}
				}
				for eventType, sequenceExecutions := range service.ExecutedSequencesPerType {
					target.IncreaseExecutedSequenceCountForType(projectName, serviceName, eventType, sequenceExecutions)
				}
			}
		}
	}
	return target
}
