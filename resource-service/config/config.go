package config

var Global EnvConfig

type EnvConfig struct {
	LogLevel                string `envconfig:"LOG_LEVEL" default:"info"`
	DirectoryStageStructure bool   `envconfig:"DIRECTORY_STAGE_STRUCTURE" default:"true"`
}
