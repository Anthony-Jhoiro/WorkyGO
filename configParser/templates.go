package configParser

// Workflow File

type WorkflowStep interface {
}

type StepDockerPersistFormat struct {
	Name   string
	Source string
}

type StepDockerFormat struct {
	Name      string   `json:"name"`
	Image     string   `json:"image"`
	Workdir   string   `json:"workdir,omitempty"`
	Commands  string   `json:"commands"`
	DependsOn []string `json:"depends_on,omitempty"`
	//Env      *map[string]string
	Persist []StepDockerPersistFormat `json:"persist,omitempty"`
}

type StepDocker struct {
	Id        string
	Name      string
	Image     string
	Workdir   string
	Commands  string
	DependsOn []string
	Persist   []StepDockerPersistFormat
}

type StepImportedFormat struct {
	Name      string
	Workflow  string
	Env       map[string]string
	Arguments map[string]string
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
