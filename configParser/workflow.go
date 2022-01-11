package configParser

import (
	"Workflow/workflow"
	"fmt"
)

type Runner struct {
	Name        string
	Description string
	Maintainer  string
	Arguments   string
	Workflow    *workflow.Workflow
}

func BuildWorkflow(template WorkflowFileTemplate, metadata WorkflowMetadataTemplate) (*Runner, error) {

	var runner Runner

	runner.Name = metadata.Name
	runner.Description = metadata.Description
	runner.Maintainer = metadata.Maintainer

	steps := make([]*workflow.Step, 0, len(template.Workflow.Steps))

	for _, stepTemplate := range template.Workflow.Steps {
		step, err := MapDictToStep(stepTemplate)
		if err != nil {
			return nil, err
		}
		// TODO : manage other steps
		dockerStep := step.(StepDockerFormat)
		s := dockerStep.ToWorkFlowStep()

		// Manage requirements
		if len(dockerStep.DependsOn) != 0 {
			for _, dependencyName := range dockerStep.DependsOn {
				dependency := findDependency(steps, dependencyName)
				if dependency == nil {
					return nil, fmt.Errorf("fail to find dependency %s", dependencyName)
				}

				s.AddRequirement(dependency)
			}
		} else {
			// TODO : Make it depends on a placeholder step
		}

		steps = append(steps, s)

	}
	runner.Workflow = workflow.NewWorkflow(steps[0], steps)

	return &runner, nil
}

func findDependency(dependencies []*workflow.Step, dependencyName string) *workflow.Step {
	for _, dependency := range dependencies {
		if dependency.GetName() == dependencyName {
			return dependency
		}
	}
	return nil
}
