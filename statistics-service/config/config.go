package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// EnvConfig godoc
type EnvConfig struct {
	AggregationIntervalSeconds int    `envconfig:"AGGREGATION_INTERVAL_SECONDS" default:"1800"`
	NextGenEvents              bool   `envconfig:"NEXT_GEN_EVENTS" default:"false"`
	LogLevel                   string `envconfig:"LOG_LEVEL" default:"info"`
	DataMigrationDisabled      bool   `envconfig:"DATA_MIGRATION_DISABLED" default:"false"`
	DataMigrationBatchSize     int    `envconfig:"DATA_MIGRATION_BATCH_SIZE" default:"150"`
	DataMigrationIntervalSec   int64  `envconfig:"DATA_MIGRATION_INTERVAL_SECONDS" default:"30"`
}

var env EnvConfig

// GetConfig godoc
func GetConfig() EnvConfig {

	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	return env
}
