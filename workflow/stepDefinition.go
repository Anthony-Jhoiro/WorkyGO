package workflow

type StepDefinition interface {
	Run(*interface{}) error
	GetLabel() string
	GetName() string
	GetDescription() string
}
