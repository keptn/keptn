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
	return [...]string{"all", "quality-gates"}[i]
}

var usecaseToID = map[string]usecase{
	"all":           AllUseCases,
	"quality-gates": QualityGates,
}
