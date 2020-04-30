package cmd

// ingressInstallOption is an enum for the Ingress installation
type ingressInstallOption int

const (
	// StopIfInstalled stops the Keptn installation if the Ingress is already installed
	StopIfInstalled ingressInstallOption = iota
	// Reuse reuses the available Ingress installation
	Reuse
	// Overwrite overwrites the Ingress installation
	Overwrite
)

func (i ingressInstallOption) String() string {
	return [...]string{"StopIfInstalled", "Reuse", "Overwrite"}[i]
}

var ingressInstallOptionToID = map[string]ingressInstallOption{
	"StopIfInstalled": StopIfInstalled,
	"Reuse":           Reuse,
	"Overwrite":       Overwrite,
}
