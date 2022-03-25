package config

// EnvConfig holds the parsed environment variables
// TODO: add other environment variables supported by Shippy
type EnvConfig struct {
	// ProjectNameMaxSize is the maximum number of characters
	// a Keptn project is allowed to have
	ProjectNameMaxSize int `envconfig:"PROJECT_NAME_MAX_SIZE" default:"200"`
}
