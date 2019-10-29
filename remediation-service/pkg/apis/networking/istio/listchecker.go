package istio

// ListChecker is an Istio type
type ListChecker struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Spec struct {
		CompiledAdapter string `json:"compiledAdapter"`
		Params          struct {
			Overrides []string `json:"overrides"`
			Blacklist bool     `json:"blacklist"`
			EntryType string   `json:"entryType"`
		} `json:"params"`
	} `json:"spec"`
}
