package backend

import "sort"

func (k K8sSecretBackend) GetScopes() ([]string, error) {
	scopes, err := k.ScopesRepository.Read()
	if err != nil {
		return nil, err
	}
	scopeArray := make([]string, len(scopes.Scopes))

	i := 0
	for scope := range scopes.Scopes {
		scopeArray[i] = scope
		i++
	}
	sort.Strings(scopeArray)
	return scopeArray, nil
}
