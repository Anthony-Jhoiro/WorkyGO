package stepMapper

import (
	"Workflow/configParser"
	"Workflow/workflow"
	"Workflow/workflow/ctx"
	"encoding/json"
	"fmt"
	"strings"
)

type StepWorkflow struct {
	Id           string            `json:"id,omitempty"`
	Name         string            `json:"name"`
	Workflow     string            `json:"workflow"`
	Parameters   map[string]string `json:"parameters,omitempty"`
	DependsOn    []string          `json:"depends_on,omitempty"`
	wf           *workflow.Workflow
	innerContext ctx.WorkflowContext
}

func (ws *StepWorkflow) GetDependencies() []string {
	if ws.DependsOn == nil {
		return []string{}
	}
	return ws.DependsOn
}

func MapWorkflowStep(template interface{}) (*StepWorkflow, error) {

	jsonBytes, err := json.Marshal(template)
	if err != nil {
		return nil, fmt.Errorf("fail to parse step : %v", err)
	}

	var workflowStep StepWorkflow

	err = json.Unmarshal(jsonBytes, &workflowStep)
	if err != nil {
		return nil, fmt.Errorf("invalid workflow step")
	}
	return &workflowStep, nil
}

func (ws *StepWorkflow) Init(ctx ctx.WorkflowContext) error {
	fileContent, err := ctx.GetExternalTemplate(ws.Workflow)

	parsedWorkflow, err := configParser.ParseWorkflowFile(fileContent, ws.Parameters)
	if err != nil {
		return fmt.Errorf("fail to parse workflow file %v", err)
	}

	wf, err := ParseWorkflowSteps(*parsedWorkflow)

	if err != nil {
		return fmt.Errorf("fail to parse workflow steps %v", err)
	}

	parentLogger := ctx.GetLogger()
	stepLogger := parentLogger.Copy("-")
	parsedWorkflow.SetLogger(stepLogger)

	ws.wf = wf
	ws.innerContext = parsedWorkflow

	return nil
}

func (ws *StepWorkflow) Clean() {

}

func (ws *StepWorkflow) GetLabel() string {
	return strings.ToLower(strings.ReplaceAll(ws.Name, " ", "_"))
}

func (ws *StepWorkflow) GetName() string {
	return ws.Name
}

func (ws *StepWorkflow) GetDescription() string {
	return ""
}

func (ws *StepWorkflow) Run(_ ctx.WorkflowContext) error {

	ws.wf.Run(ws.innerContext)
	return nil
}
