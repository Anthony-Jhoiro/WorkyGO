package configParser

// This file contains all struct definitions to parse the yaml template

type workflowWorkflowFormat struct {
	Steps []interface{} `yaml:"steps"`
}

type workflowFileTemplate struct {
	Workflow workflowWorkflowFormat `yaml:"workflow"`
}

// Config file

type parameterTemplate struct {
	Name         string                `yaml:"name"`
	Description  string                `yaml:"description,omitempty"`
	Validators   string                `yaml:"validators,omitempty"`
	DefaultValue string                `yaml:"default_value" yaml:"default_value"`
	Type         WorkflowParameterType `yaml:"type"`
}

// WorkflowParameterType describe the parameter types that can be used in a workflow template
type WorkflowParameterType string

const (
	NumberType  WorkflowParameterType = "number"
	BooleanType                       = "boolean"
	StringType                        = "string"
	FloatType                         = "float"
)

// importTemplate is the template format that defines an import statement
type importTemplate struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type workflowMetadataTemplate struct {
	Name        string                       `yaml:"name"`
	Description string                       `yaml:"description"`
	Maintainer  string                       `yaml:"maintainer"`
	Parameters  map[string]parameterTemplate `yaml:"parameters"`
	Imports     []importTemplate             `yaml:"imports"`
}
