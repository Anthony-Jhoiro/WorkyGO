package stepMapper

import (
	"Workflow/configParser"
	"Workflow/workflow"
	"fmt"
)

func ParseWorkflowSteps(parsedWorkflow configParser.ParsedWorkflow) (*configParser.Runner, error) {
	steps, err := mapMultipleSteps(parsedWorkflow.Steps)

	if err != nil {
		return nil, fmt.Errorf("fail to map steps : %v", err)
	}

	return buildWorkflow(parsedWorkflow, steps)
}

func buildWorkflow(parsedWorkflow configParser.ParsedWorkflow, stepsDefinitions []workflow.StepDefinition) (*configParser.Runner, error) {
	steps := make([]*workflow.Step, 0, len(stepsDefinitions))

	stepMapper := map[string]*workflow.Step{}

	for _, stepDefinition := range stepsDefinitions {
		step := &workflow.Step{
			StepDefinition: stepDefinition,
			Status:         0,
			NextSteps:      make([]*workflow.Step, 0, len(stepsDefinitions)),
			PreviousSteps:  make([]*workflow.Step, 0, len(stepDefinition.GetDependencies())),
		}
		steps = append(steps, step)
		stepMapper[step.GetLabel()] = step
	}

	initialStepDef := MakeEmptyStep()

	initialStep := &workflow.Step{
		StepDefinition: initialStepDef,
		Status:         0,
		NextSteps:      make([]*workflow.Step, 0, len(stepsDefinitions)),
		PreviousSteps:  make([]*workflow.Step, 0),
	}

	for _, step := range steps {
		deps := step.GetDependencies()
		if len(deps) == 0 {
			step.AddRequirement(initialStep)
		} else {
			for _, dep := range deps {
				if val, ok := stepMapper[dep]; ok {
					step.AddRequirement(val)
				} else {
					return nil, fmt.Errorf("step %s does not exists", dep)
				}
			}
		}
	}

	wf := workflow.NewWorkflow(initialStep, steps)

	runner := &configParser.Runner{
		Name:        parsedWorkflow.Name,
		Description: parsedWorkflow.Description,
		Maintainer:  parsedWorkflow.Maintainer,
		Arguments:   parsedWorkflow.Arguments,
		Workflow:    wf,
	}

	return runner, nil
}

func mapMultipleSteps(stepTemplates []interface{}) ([]workflow.StepDefinition, error) {
	steps := make([]workflow.StepDefinition, 0, len(stepTemplates))

	for stepIndex, stepTemplate := range stepTemplates {
		step, err := mapStep(stepTemplate)
		if err != nil {
			return nil, fmt.Errorf("step %d is invalid : %v", stepIndex, err)
		}

		steps = append(steps, step)
	}
	return steps, nil
}

func mapStep(tpl interface{}) (workflow.StepDefinition, error) {
	return MapDockerStep(tpl)
}
