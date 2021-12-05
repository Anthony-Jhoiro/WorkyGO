package configParser

// Workflow File

type WorkflowStep interface {
}

type StepDockerFormat struct {
	Name     string
	Image    string
	Workdir  string
	Commands string
	//Env      *map[string]string
	Persist []interface{}
}

type StepImportedFormat struct {
	Name       string
	Workflow   string
	Env        map[string]string
	Parameters map[string]string
}

type WorkflowWorkflowFormat struct {
	Steps []map[string]interface{}
}

type WorkflowFormat struct {
	Workflow WorkflowWorkflowFormat
}

type WorkflowFileTemplate struct {
	Workflow WorkflowWorkflowFormat
}

// Config file

type ParameterTemplate struct {
	Name         string
	Description  string
	Validators   string
	DefaultValue string `yaml:"default_value"`
	Type         WorkflowParameterType
}

type WorkflowParameterType string

//boolean, string, character, integer, floating-point
const (
	NumberType  WorkflowParameterType = "number"
	BooleanType                       = "boolean"
	StringType                        = "string"
	FloatType                         = "float"
)

type ConfigTemplate struct {
	Name       string
	Parameters []ParameterTemplate
}
