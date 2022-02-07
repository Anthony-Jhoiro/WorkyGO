package workflow

import (
	"Workflow/workflow/ctx"
	"fmt"
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

func (wf *Workflow) Run(ctx ctx.WorkflowContext) error {
	channel := make(chan *Step, wf.nodeCount)
	runningSteps := 1

	stepsOutput := make(map[string]map[string]string)

	errStack := make([]string, 0, len(wf.Steps))

	go wf.firstStep.Execute(channel, ctx, stepsOutput)

	for runningSteps != 0 {
		closingStep := <-channel
		runningSteps -= 1
		stepsOutput[closingStep.GetLabel()] = closingStep.GetOutput()

		if closingStep.Status == StepFail {
			errStack = append(errStack, closingStep.GetLabel())
		}

		// Mark following steps as failed
		for _, e := range closingStep.NextSteps {

			requirementsOk, err := e.RequirementsFulfilled()
			runningSteps += 1

			if err != nil {
				go func(executable *Step) {
					executable.fail(channel)
				}(e)
			} else if requirementsOk {
				go func(executable *Step) {
					executable.Execute(channel, ctx, stepsOutput)
				}(e)
			}
		}
	}

	if len(errStack) != 0 {
		return fmt.Errorf("some step failed : %v", errStack)
	}
	return nil

}
