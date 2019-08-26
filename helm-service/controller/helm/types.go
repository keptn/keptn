package helm

// Chart represents a Helm chart
type Chart struct {
	APIVersion  string `json:"apiVersion"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Version     string `json:"version"`
}

// Requirements represents the Helm umbrella requirements
type Requirements struct {
	Dependencies []struct {
		Name      string `json:"name"`
		Version   string `json:"version"`
		Condition string `json:"condition"`
	} `json:"dependencies"`
}

// Values represents the Helm umbrella values
type Values map[string]Enabler

// Enabler
type Enabler struct {
	Enabled bool `json:"enabled"`
}
