package helm

// GetGatewayName returns the name of the gateway for a specific project and stage
func GetGatewayName(project string, stage string) string {
	return project + "-" + stage + "-gateway"
}

// GetUmbrellaReleaseName returns the release name of the umbrella chart
func GetUmbrellaReleaseName(project string, stage string) string {
	return project + "-" + stage
}

// GetUmbrellaNamespace returns the namespace in which the umbrella chart (e.g. containing the gateway) is applied
func GetUmbrellaNamespace(project string, stage string) string {
	return project + "-" + stage
}

// GetChartName returns the name of the chart
func GetChartName(service string, generated bool) string {
	suffix := ""
	if generated {
		suffix = "-generated"
	}
	return service + suffix
}

// GetReleaseName returns the name of the Helm release
func GetReleaseName(project string, stage string, service string, generated bool) string {
	suffix := ""
	if generated {
		suffix = "-generated"
	}
	return project + "-" + stage + "-" + service + suffix
}
