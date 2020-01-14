package cmd

// istioInstallOption is an enum for the Istio installation
type istioInstallOption int

const (
	// StopIfInstalled stops the Keptn installation if Istio is already installed
	StopIfInstalled = iota
	// Reuse reuses the available Istio installation
	Reuse
	// Overwrite overwrites the Istio installation
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
