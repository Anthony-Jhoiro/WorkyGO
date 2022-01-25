package workflow

type StepDefinition interface {
	Init() error
	Run(*interface{}) error
	Clean()
	GetLabel() string
	GetName() string
	GetDescription() string
	GetDependencies() []string
}
