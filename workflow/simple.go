package workflow

type RunFunction func(*interface{}) error

type SimpleStep struct {
	label       string
	name        string
	description string
	runFunction RunFunction
}

func (step *SimpleStep) Init() error {
	return nil
}

func (step *SimpleStep) Clean() {

}

func (step *SimpleStep) GetDependencies() []string {
	return nil
}

func (step *SimpleStep) GetLabel() string {
	return step.label
}

func (step *SimpleStep) GetName() string {
	return step.name
}

func (step *SimpleStep) GetDescription() string {
	return step.description
}

func (step *SimpleStep) Run(ctx *interface{}) error {
	return step.runFunction(ctx)
}

func NewSimpleStep(label string, name string, description string, runFunction RunFunction) *Step {
	return &Step{
		StepDefinition: &SimpleStep{
			label:       label,
			name:        name,
			description: description,
			runFunction: runFunction,
		},
		Status:        StepPending,
		NextSteps:     nil,
		PreviousSteps: nil,
	}
}
