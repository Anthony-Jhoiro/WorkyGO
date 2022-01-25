package configParser

// Workflow File

type workflowWorkflowFormat struct {
	Steps []interface{}
}

type workflowFileTemplate struct {
	Workflow workflowWorkflowFormat
}

// Config file

type parameterTemplate struct {
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

type importTemplate struct {
	Name string
	Url  string
}

type workflowMetadataTemplate struct {
	Name        string                       `yaml:"name"`
	Description string                       `yaml:"description"`
	Maintainer  string                       `yaml:"maintainer"`
	Parameters  map[string]parameterTemplate `yaml:"parameters"`
	Imports     []importTemplate             `yaml:"imports"`
}
