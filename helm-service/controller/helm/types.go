package helm

// Requirements represents the Helm umbrella requirements
type Requirements struct {
	Dependencies []RequirementDependencies `json:"dependencies"`
}

// RequirementDependencies represents the dependencies contained in the Helm umbrella requirements
type RequirementDependencies struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Condition  string `json:"condition"`
	Repository string `json:"repository"`
}

// Values represents the Helm umbrella values
type Values map[string]Enabler

// Enabler in Values file
type Enabler struct {
	Enabled bool `json:"enabled"`
}
