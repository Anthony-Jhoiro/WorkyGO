package steps

import (
	"container/list"
	"log"
)

func (step *SimpleStep) Execute(ctx interface{}) (interface{}, error) {
	// Check if requirements are fulfilled

	for e := step.Requirements.Front(); e != nil; e = e.Next() {
		if e.Value.(ExecutableStep).GetStatus() != StepOK {
			return nil, nil
		}
	}

	// Execute the steps
	step.Status = StepRunning

	// Replace with container execution
	log.Printf("-> %s", step.Label)

	step.Status = StepOK

	for e := step.NextSteps.Front(); e != nil; e = e.Next() {

		_, err := e.Value.(ExecutableStep).Execute(ctx)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (step *SimpleStep) GetStatus() StepStatus {
	return step.Status
}

func CreateStep(name string) *SimpleStep {
	return &SimpleStep{
		Step: Step{
			Name:         name,
			Description:  "lorem ipsum",
			Label:        name,
			Status:       StepPending,
			NextSteps:    list.New(),
			Requirements: list.New(),
		},
		DockerfileUrl:   "",
		DockerImageName: "",
	}
}

func (step *SimpleStep) AddRequirement(parent *SimpleStep) {
	step.Requirements.PushBack(parent)
	parent.NextSteps.PushBack(step)
}

func (s *SimpleStep) Print() {
	log.Printf("Node %s", s.Name)
}
