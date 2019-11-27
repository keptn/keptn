package cmd

// Ingress is an enum type for the ingress
type Ingress int

const (
	Istio Ingress = iota
	Nginx
)

func (i Ingress) String() string {
	return [...]string{"istio", "nginx"}[i]
}
