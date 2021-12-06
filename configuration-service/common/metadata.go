package common

type ProjectMetadata struct {
	ProjectName               string `yaml:"projectName"`
	CreationTimestamp         string `yaml:"creationTimestamp"`
	IsUsingDirectoryStructure bool   `yaml:"isUsingDirectoryStructure"`
}

type StageMetadata struct {
	StageName         string
	CreationTimestamp string
}

type ServiceMetadata struct {
	ServiceName       string
	CreationTimestamp string
}
