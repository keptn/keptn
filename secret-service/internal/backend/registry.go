package backend

var backendRegistry = map[string]func() SecretBackend{}

func Register(name string, factory func() SecretBackend) {
	backendRegistry[name] = factory
}

func GetRegisteredBackends() []string {
	var r []string
	for i, _ := range backendRegistry {
		r = append(r, i)
	}
	return r
}

func CreateBackend(backendType string) SecretBackend {
	return backendRegistry[backendType]()
}
