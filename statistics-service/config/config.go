package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

// EnvConfig godoc
type EnvConfig struct {
	AggregationIntervalSeconds int  `envconfig:"AGGREGATION_INTERVAL_SECONDS" default:"1800"`
	NextGenEvents              bool `envconfig:"NEXT_GEN_EVENTS" default:"false"`
}

var env EnvConfig

// GetConfig godoc
func GetConfig() EnvConfig {

	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	return env
}
