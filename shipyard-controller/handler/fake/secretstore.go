package fake

type MockSecretStore struct {
	CreateFunc func(name string, content map[string][]byte) error
	DeleteFunc func(name string) error
	GetFunc    func(name string) (map[string][]byte, error)
	UpdateFunc func(name string, content map[string][]byte) error
}

func (ms *MockSecretStore) CreateSecret(name string, content map[string][]byte) error {
	return ms.CreateFunc(name, content)
}

func (ms *MockSecretStore) DeleteSecret(name string) error {
	return ms.DeleteFunc(name)
}

func (ms *MockSecretStore) GetSecret(name string) (map[string][]byte, error) {
	return ms.GetFunc(name)
}

func (ms *MockSecretStore) UpdateSecret(name string, content map[string][]byte) error {
	return ms.UpdateFunc(name, content)
}
