package cmd

// istioInstallOption is an enum for the Istio installation
type istioInstallOption int

const (
	// StopIfAvailable stops the Keptn installation if Istio is available
	StopIfAvailable = iota
	// Reuse reuses the available Istio installation
	Reuse
	// Overwrite overwrites the Istio installation
	Overwrite
)

func (i istioInstallOption) String() string {
	return [...]string{"StopIfAvailable", "Reuse", "Overwrite"}[i]
}

var istioInstallOptionToID = map[string]istioInstallOption{
	"StopIfAvailable": StopIfAvailable,
	"Reuse":           Reuse,
	"Overwrite":       Overwrite,
}
