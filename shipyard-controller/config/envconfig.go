package config

// EnvConfig holds the parsed environment variables
// TODO: add other environment variables supported by Shippy
type EnvConfig struct {
	// ProjectNameMaxLength is the maximum number of characters
	// a Keptn project is allowed to have
	ProjectNameMaxLength int `envconfig:"PROJ_NAME_MAX_LENGTH" default:"100"`
}
