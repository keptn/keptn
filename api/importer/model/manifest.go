package model

type APITask struct {
	Action      string `yaml:"action"`
	PayloadFile string `yaml:"payload"`
}

type ResourceTask struct {
	File      string `yaml:"resource"`
	RemoteURI string `yaml:"resourceUri"`
	Stage     string `yaml:"stage"`
}

type ManifestTask struct {
	*APITask      `yaml:",inline"`
	*ResourceTask `yaml:",inline"`
	ID            string `yaml:"id"`
	Type          string `yaml:"type"`
	Name          string `yaml:"name"`
}

type ImportManifest struct {
	ApiVersion string          `yaml:"apiVersion"`
	Tasks      []*ManifestTask `yaml:"tasks"`
}
