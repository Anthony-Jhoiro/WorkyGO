package stepMapper

import (
	"encoding/json"
	"fmt"
	"github.com/Anthony-Jhoiro/WorkyGO/configParser"
	"github.com/Anthony-Jhoiro/WorkyGO/workflow"
	"github.com/Anthony-Jhoiro/WorkyGO/workflow/ctx"
	"strings"
	"text/template"
)

type StepWorkflow struct {
	Id           string            `json:"id,omitempty"`
	Name         string            `json:"name"`
	Workflow     string            `json:"workflow"`
	Parameters   map[string]string `json:"parameters,omitempty"`
	DependsOn    []string          `json:"depends_on,omitempty"`
	outputs      map[string]string
	wf           *workflow.Workflow
	innerContext ctx.WorkflowContext
	outputValues map[string]string
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

func (ws *StepWorkflow) resolveStepValues(previousStepsOutput map[string]map[string]string) error {
	parser := template.New("parser").Funcs(template.FuncMap{"getVar": getOutputParser(previousStepsOutput)})

	for key, parameter := range ws.Parameters {
		err := _resolveStepValue(parser, &parameter)
		if err != nil {
			return fmt.Errorf("fail to parse parameter %s : %v", key, err)
		}
	}

	return nil
}

func (ws *StepWorkflow) Init(ctx ctx.WorkflowContext, previousStepsOutput map[string]map[string]string) error {

	err := ws.resolveStepValues(previousStepsOutput)
	if err != nil {
		return fmt.Errorf("fail to parse some values : %v", err)
	}

	fileContent, err := ctx.GetExternalTemplate(ws.Workflow)

	parsedWorkflow, err := configParser.ParseWorkflowFile(fileContent, ws.Parameters)
	if err != nil {
		return fmt.Errorf("fail to parse workflow file %v", err)
	}

	wf, err := ParseWorkflowSteps(*parsedWorkflow)

	ws.outputs = parsedWorkflow.Output

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
	output, err := ws.wf.Run(ws.innerContext)

	if err != nil {
		return err
	}

	parser := getOutputParser(output)

	ov := make(map[string]string)

	for outputName, key := range ws.outputs {
		slicedKey := strings.Split(key, ".")

		if len(slicedKey) != 2 {
			return fmt.Errorf("output key %s is malformed", key)
		}
		value, err := parser(slicedKey[0], slicedKey[1])
		if err != nil {
			return fmt.Errorf("no value founded for key %s : %v", key, err)

		}
		ov[outputName] = value
	}

	ws.outputValues = ov
	return nil
}

func (ws *StepWorkflow) GetOutput() map[string]string {
	return ws.outputValues
}
