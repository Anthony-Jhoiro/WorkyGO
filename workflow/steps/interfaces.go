package steps

import (
	"container/list"
	"log"
)

type StepStatus uint16

const (
	StepOK      StepStatus = 0
	StepFail               = 1
	StepRunning            = 2
	StepPending            = 3
)

type StepConfig struct {
	Name               string
	Description        string
	Label              string
	DockerfileLocation string
	DockerImageName    string
}

type Workflow struct {
}

type ExecutableStep interface {
	Execute(interface{}) (interface{}, error)
	GetStatus() StepStatus
}

type Step struct {
	Name         string
	Description  string
	Label        string
	Status       StepStatus
	Requirements *list.List
	NextSteps    *list.List
}

type SimpleStep struct {
	Step
	DockerfileUrl   string
	DockerImageName string
}

type ComplexStep struct {
	Step
	FirstStep ExecutableStep
}

func (s *ComplexStep) Execute(_ctx interface{}) (interface{}, error) {

	log.Printf("Executng step '%s'", s.Name)

	return s.FirstStep.Execute(_ctx)
}
