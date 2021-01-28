package fake

type SecretStore struct {
	CreateFunc func(name string, content map[string][]byte) error
	DeleteFunc func(name string) error
	GetFunc    func(name string) (map[string][]byte, error)
	UpdateFunc func(name string, content map[string][]byte) error
}

func (ms *SecretStore) CreateSecret(name string, content map[string][]byte) error {
	return ms.CreateFunc(name, content)
}

func (ms *SecretStore) DeleteSecret(name string) error {
	return ms.DeleteFunc(name)
}

func (ms *SecretStore) GetSecret(name string) (map[string][]byte, error) {
	return ms.GetFunc(name)
}

func (ms *SecretStore) UpdateSecret(name string, content map[string][]byte) error {
	return ms.UpdateFunc(name, content)
}
