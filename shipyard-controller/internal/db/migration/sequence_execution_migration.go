package migration

import (
	"fmt"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	v1 "github.com/keptn/keptn/shipyard-controller/internal/db/models/sequence_execution/v1"
	"github.com/keptn/keptn/shipyard-controller/models"
	logger "github.com/sirupsen/logrus"
)

// NewSequenceExecutionMigrator creates a new SequenceExecutionMigrator
// Internally it is using the SequenceExecutionJsonStringRepo decorator
// which stores the arbitrary event payload sent by keptn integrations as Json strings to avoid having property names with dots (.) in them
func NewSequenceExecutionMigrator(dbConnection *db.MongoDBConnection) *SequenceExecutionMigrator {
	return &SequenceExecutionMigrator{
		projectRepo:           db.NewMongoDBKeyEncodingProjectsRepo(dbConnection),
		sequenceExecutionRepo: db.NewMongoDBSequenceExecutionRepo(dbConnection, db.WithSequenceExecutionModelTransformer(&v1.ModelTransformer{})),
	}
}

type SequenceExecutionMigrator struct {
	sequenceExecutionRepo db.SequenceExecutionRepo
	projectRepo           db.ProjectRepo
}

// Run retrieves all existing sequence executions from the repository
// and performs an update operation on each of them using the SequenceExecutionJsonStringRepo.
// This way, sequence executions containing stored with the previous format are migrated to the new one
func (s *SequenceExecutionMigrator) Run() error {
	projects, err := s.projectRepo.GetProjects()
	if err != nil {
		return fmt.Errorf("could not migrate sequence executions to new format: %w", err)
	}
	s.updateSequenceExecutionsOfProject(projects)
	return nil
}

func (s *SequenceExecutionMigrator) updateSequenceExecutionsOfProject(projects []*apimodels.ExpandedProject) {
	if projects == nil {
		return
	}
	for _, project := range projects {
		sequenceExecutions, err := s.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{
			Scope: models.EventScope{
				EventData: keptnv2.EventData{
					Project: project.ProjectName,
				},
			}})
		if err != nil {
			logger.Errorf("Could not retrieve sequence executions for project %s: %v", project.ProjectName, err)
			continue
		}
		for _, sequenceExecution := range sequenceExecutions {
			if sequenceExecution.SchemaVersion == v1.SchemaVersionV1 {
				continue
			}
			if err := s.sequenceExecutionRepo.Upsert(sequenceExecution, &models.SequenceExecutionUpsertOptions{
				Replace: true,
			}); err != nil {
				logger.Errorf("Could not update sequence execution with ID %s for project %s: %v", sequenceExecution.ID, project.ProjectName, err)
				continue
			}
		}
	}
}
