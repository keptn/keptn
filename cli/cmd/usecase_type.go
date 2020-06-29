package cmd

// usecase is an enum for the supported use cases
type usecase int

const (
	// All supports all Keptn use cases
	AllUseCases usecase = iota
	// QualityGates supports only quality gates use cases
	QualityGates
)

func (i usecase) String() string {
	return [...]string{"continuous-delivery", "quality-gates"}[i]
}

var usecaseToID = map[string]usecase{
	"continuous-delivery": AllUseCases,
	"quality-gates":       QualityGates,
}
