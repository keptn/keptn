package backend

var backendRegistry = map[string]func() SecretBackend{}

func Register(name string, factory func() SecretBackend) {
	backendRegistry[name] = factory
}

func GetRegisteredBackends() []string {
	r := make([]string, len(backendRegistry))
	for i := range backendRegistry {
		r = append(r, i)
	}
	return r
}

func CreateBackend(backendType string) SecretBackend {
	return backendRegistry[backendType]()
}
