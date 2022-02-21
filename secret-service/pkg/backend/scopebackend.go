package backend

type ScopeBackend interface {
	GetScopes() ([]string, error)
}
