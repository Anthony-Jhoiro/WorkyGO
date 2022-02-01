package stepMapper

import (
	"Workflow/configParser"
	"Workflow/workflow"
	"fmt"
	"strings"
)

func ParseWorkflowSteps(parsedWorkflow configParser.ParsedWorkflow) (*workflow.Workflow, error) {
	steps, err := mapMultipleSteps(parsedWorkflow.Steps)

	if err != nil {
		return nil, fmt.Errorf("fail to map steps : %v", err)
	}

	return buildWorkflow(steps)
}

func buildWorkflow(stepsDefinitions []workflow.StepDefinition) (*workflow.Workflow, error) {
	steps := make([]*workflow.Step, 0, len(stepsDefinitions)+1)

	stepMapper := map[string]*workflow.Step{}

	for _, stepDefinition := range stepsDefinitions {
		step := &workflow.Step{
			StepDefinition: stepDefinition,
			Status:         workflow.StepPending,
			NextSteps:      make([]*workflow.Step, 0, len(stepsDefinitions)),
			PreviousSteps:  make([]*workflow.Step, 0, len(stepDefinition.GetDependencies())),
		}
		steps = append(steps, step)
		stepMapper[step.GetLabel()] = step
	}

	initialStepDef := MakeEmptyStep()

	initialStep := &workflow.Step{
		StepDefinition: initialStepDef,
		Status:         workflow.StepPending,
		NextSteps:      make([]*workflow.Step, 0, len(stepsDefinitions)),
		PreviousSteps:  make([]*workflow.Step, 0),
	}

	for _, step := range steps {
		deps := step.GetDependencies()
		if len(deps) == 0 {
			step.AddRequirement(initialStep)
		} else {
			for _, dep := range deps {
				labelDep := strings.ToLower(strings.ReplaceAll(dep, " ", "_"))
				if val, ok := stepMapper[labelDep]; ok {
					step.AddRequirement(val)
				} else {
					return nil, fmt.Errorf("step %s does not exists", labelDep)
				}
			}
		}
	}

	steps = append(steps, initialStep)

	wf := workflow.NewWorkflow(initialStep, steps)

	return wf, nil
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
	var step workflow.StepDefinition

	step, err := MapDockerStep(tpl)

	if err == nil {
		return step, nil
	}

	step, err = MapWorkflowStep(tpl)

	if err == nil {
		return step, nil
	}

	return nil, fmt.Errorf("invalid step")
}
