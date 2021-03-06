package workflow

import (
	"fmt"
	"github.com/Anthony-Jhoiro/WorkyGO/workflow/ctx"
)

type Workflow struct {
	nodeCount int
	firstStep *Step
	Steps     []*Step
}

func NewWorkflow(firstStep *Step, steps []*Step) *Workflow {
	return &Workflow{
		nodeCount: len(steps),
		firstStep: firstStep,
		Steps:     steps,
	}
}

func (wf *Workflow) Print() {
	for _, step := range wf.Steps {
		fmt.Printf("[ %s ] => %s \n", step.GetLabel(), step.Status.GetName())
	}
}

func (wf *Workflow) Run(ctx ctx.WorkflowContext) (map[string]map[string]string, error) {
	channel := make(chan *Step, wf.nodeCount)
	runningSteps := new(int32)
	*runningSteps = 1

	stepsOutput := make(map[string]map[string]string)

	errStack := make([]string, 0, len(wf.Steps))

	go wf.firstStep.Execute(channel, ctx, stepsOutput)

	for *runningSteps != 0 {
		closingStep := <-channel
		//*runningSteps -= 1
		*runningSteps -= 1
		stepsOutput[closingStep.GetLabel()] = closingStep.GetOutput()

		if closingStep.Status == StepFail {
			errStack = append(errStack, closingStep.GetLabel())
		}

		// Mark following steps as failed
		for _, e := range closingStep.NextSteps {

			requirementsOk, err := e.RequirementsFulfilled()

			if err != nil {
				*runningSteps += 1
				go func(executable *Step) {
					executable.fail(channel)
				}(e)
			} else if requirementsOk {
				*runningSteps += 1
				go func(executable *Step) {
					executable.Execute(channel, ctx, stepsOutput)
				}(e)
			}
		}
	}

	if len(errStack) != 0 {
		return nil, fmt.Errorf("some step failed : %v", errStack)
	}
	return stepsOutput, nil

}
