package cmd

// Ingress is an enum type for the ingress
type Ingress int

const (
	istio Ingress = iota
	nginx
)

func (i Ingress) String() string {
	return [...]string{"istio", "nginx"}[i]
}
