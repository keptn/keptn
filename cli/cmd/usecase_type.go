package cmd

// usecase is an enum for the supported use cases
type usecase int

const (
	// QualityGates supports only quality gates use cases
	QualityGates usecase = iota

	// ContinuousDelivery supports all Keptn use cases
	ContinuousDelivery
)

func (i usecase) String() string {
	return [...]string{"quality-gates", "continuous-delivery"}[i]
}

var usecaseToID = map[string]usecase{
	"continuous-delivery": ContinuousDelivery,
	"quality-gates":       QualityGates,
	"":                    QualityGates,
}
