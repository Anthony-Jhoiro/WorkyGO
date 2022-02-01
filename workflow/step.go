package workflow

import "Workflow/workflow/ctx"

type StepStatus uint16

type Step struct {
	StepDefinition
	Status        StepStatus
	NextSteps     []*Step
	PreviousSteps []*Step
}

// RequirementsFulfilled Check if all requirements steps have been completed by checking there status.
func (step *Step) RequirementsFulfilled() bool {
	for _, requirement := range step.PreviousSteps {
		if requirement.Status != StepOK {
			return false
		}
	}
	return true
}

func (step *Step) AddRequirement(parent *Step) {
	parent.NextSteps = append(parent.NextSteps, step)
	step.PreviousSteps = append(step.PreviousSteps, parent)
}

func (step *Step) Execute(channel chan *Step, ctx ctx.WorkflowContext) {
	// Execute the steps
	step.Status = StepRunning

	err := step.Init(ctx)
	if err != nil {
		step.Status = StepFail
		channel <- step
		return
	}

	if step.Run(ctx) != nil {
		step.Status = StepFail
	} else {
		step.Status = StepOK
	}

	step.Clean()

	channel <- step
}
