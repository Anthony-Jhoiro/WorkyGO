package stepMapper

import "Workflow/workflow"

type EmptyStep struct {
}

func MakeEmptyStep() workflow.StepDefinition {
	return &EmptyStep{}
}

func (es *EmptyStep) GetDependencies() []string {
	panic("implement me")
}

func (es *EmptyStep) Clean() {

}

func (es *EmptyStep) GetLabel() string {
	return "empty"
}

func (es *EmptyStep) GetName() string {
	return "Empty Step"
}

func (es *EmptyStep) GetDescription() string {
	return "This is an empty step used to start a workflow."
}

func (es *EmptyStep) Init() error {
	return nil
}

func (es *EmptyStep) Run(_ *interface{}) error {
	return nil
}
