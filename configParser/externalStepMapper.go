package configParser

import (
	"Workflow/workflow"
	"fmt"
	"os"
	"path"
)

const BASE_WORKFLOW_DIRECTORY = "./workflows"

type ExternalStep struct {
	Id           string
	Name         string
	DependsOn    []string
	WorkflowName string
	Arguments    map[string]string
}

func (step *ExternalStep) ToWorkFlowStep() *workflow.Step {
	return workflow.NewSimpleStep(step.Name, step.Name, "", step.Run)
}

func (step ExternalStep) GetFileName() string {
	return path.Join(BASE_WORKFLOW_DIRECTORY, fmt.Sprintf("%s.yaml", step.WorkflowName))
}

func (step *ExternalStep) TemplateFileExists() bool {
	_, err := os.Stat(step.GetFileName())
	return os.IsNotExist(err)
}

func (step *ExternalStep) Run(_ *interface{}) error {
	// Search for file on system
	templateFile, err := os.ReadFile(step.GetFileName())
	if err != nil {
		return fmt.Errorf("fail to read external workflow %s : %v", step.WorkflowName, err)
	}

	res, err := ParseWorkflowFile(templateFile, step.Arguments)
	if err != nil {
		return fmt.Errorf("fail to parse external workflow %s : %v", step.WorkflowName, err)
	}

	res.Workflow.Run()

	return nil
}
