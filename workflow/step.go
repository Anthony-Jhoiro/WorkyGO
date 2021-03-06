package workflow

import (
	"fmt"
	"github.com/Anthony-Jhoiro/WorkyGO/workflow/ctx"
)

type StepStatus uint16

type Step struct {
	StepDefinition
	Status        StepStatus
	NextSteps     []*Step
	PreviousSteps []*Step
}

// RequirementsFulfilled Check if all requirements steps have been completed by checking there status.
func (step *Step) RequirementsFulfilled() (bool, error) {
	for _, requirement := range step.PreviousSteps {
		if requirement.Status == StepFail {
			return true, fmt.Errorf("a requirement of step %s failed", step.GetLabel())
		}
		if requirement.Status != StepOK {
			return false, nil
		}
	}
	return true, nil
}

func (step *Step) AddRequirement(parent *Step) {
	parent.NextSteps = append(parent.NextSteps, step)
	step.PreviousSteps = append(step.PreviousSteps, parent)
}

func (step *Step) Execute(channel chan *Step, ctx ctx.WorkflowContext, previousStepsOutput map[string]map[string]string) {
	// The step must be written in the channel at the end of the function
	defer func() {
		channel <- step
	}()

	stepContext := ctx.Copy()
	l := ctx.GetLogger()
	sl := l.Copy(fmt.Sprintf("[%s]", step.GetLabel()))
	stepContext.SetLogger(sl)

	// Execute the Steps
	step.Status = StepRunning
	_ = stepContext.GetLogger().Print(fmt.Sprintf("Starting %s", step.GetLabel()))

	err := step.Init(stepContext, previousStepsOutput)
	if err != nil {
		step.Status = StepFail
		_ = stepContext.GetLogger().Print(fmt.Sprintf("[ERROR] Step %s initialization failed : %v", step.GetLabel(), err))
		return
	}

	if err = step.Run(stepContext); err != nil {
		step.Status = StepFail
		_ = stepContext.GetLogger().Print(fmt.Sprintf("[ERROR] Step %s failed : %v", step.GetLabel(), err))

	} else {
		step.Status = StepOK
		_ = stepContext.GetLogger().Print(fmt.Sprintf("Step %s succeeded", step.GetLabel()))
	}
}

func (step *Step) fail(channel chan *Step) {
	step.Status = StepFail
	channel <- step
}
