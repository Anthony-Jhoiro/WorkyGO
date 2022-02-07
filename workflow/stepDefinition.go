package workflow

import "Workflow/workflow/ctx"

type StepDefinition interface {
	Init(ctx ctx.WorkflowContext, previousStepsOutput map[string]map[string]string) error
	Run(ctx ctx.WorkflowContext) error
	Clean()
	GetLabel() string
	GetName() string
	GetDescription() string
	GetDependencies() []string
	GetOutput() map[string]string
}
