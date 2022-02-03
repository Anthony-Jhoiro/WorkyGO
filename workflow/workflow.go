package workflow

import (
	"Workflow/workflow/ctx"
	"fmt"
	"os"
	"os/exec"
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

func (wf *Workflow) Run(ctx ctx.WorkflowContext) {
	channel := make(chan *Step, wf.nodeCount)

	go wf.firstStep.Execute(channel, ctx)

	for i := 0; i < wf.nodeCount; i++ {
		closingStep := <-channel

		for _, e := range closingStep.NextSteps {
			if e.RequirementsFulfilled() {
				// Execute step
				go func(executable *Step) {
					executable.Execute(channel, ctx)
				}(e)
			}
		}

		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
		wf.Print()
	}
}
