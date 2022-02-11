package stepMapper

import (
	"github.com/Anthony-Jhoiro/WorkyGO/workflow"
	"github.com/Anthony-Jhoiro/WorkyGO/workflow/ctx"
)

type EmptyStep struct {
}

func MakeEmptyStep() workflow.StepDefinition {
	return &EmptyStep{}
}

func (es *EmptyStep) GetDependencies() []string {
	return []string{}
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

func (es *EmptyStep) Init(_ ctx.WorkflowContext) error {
	return nil
}

func (es *EmptyStep) Run(_ ctx.WorkflowContext) error {
	return nil
}
