package cmd

// istioInstallOption is an enum for the istio installation
type istioInstallOption int

const (
	// StopIfInstalled stops the Keptn installation if istio is already installed
	StopIfInstalled istioInstallOption = iota
	// Reuse reuses the available istio installation
	Reuse
	// Overwrite overwrites the istio installation
	Overwrite
)

func (i istioInstallOption) String() string {
	return [...]string{"StopIfInstalled", "Reuse", "Overwrite"}[i]
}

var istioInstallOptionToID = map[string]istioInstallOption{
	"StopIfInstalled": StopIfInstalled,
	"Reuse":           Reuse,
	"Overwrite":       Overwrite,
}
