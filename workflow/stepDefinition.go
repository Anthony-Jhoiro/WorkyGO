package workflow

import "github.com/Anthony-Jhoiro/WorkyGO/workflow/ctx"

type StepDefinition interface {
	Init(ctx ctx.WorkflowContext) error
	Run(ctx ctx.WorkflowContext) error
	Clean()
	GetLabel() string
	GetName() string
	GetDescription() string
	GetDependencies() []string
}
