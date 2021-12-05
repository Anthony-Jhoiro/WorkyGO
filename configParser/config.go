package configParser

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"strconv"
	"text/template"
)

func ParseConfigFileContent(fileContent []byte) (*ConfigTemplate, error) {

	tpl := &ConfigTemplate{}
	err := yaml.Unmarshal(fileContent, &tpl)

	return tpl, err
}

type ImportTemplate struct {
	Name string
	Url  string
}

type WorkflowMetadataTemplate struct {
	Name        string                       `yaml:"name"`
	Description string                       `yaml:"description"`
	Maintainer  string                       `yaml:"maintainer"`
	Parameters  map[string]ParameterTemplate `yaml:"parameters"`
	Imports     []ImportTemplate             `yaml:"imports"`
}

// ParseWorkflowFile Parse the content of a workflow file with the given arguments
func ParseWorkflowFile(fileContent []byte, arguments map[string]string) (*Runner, error) {
	metadata, err := LoadMetadata(fileContent)
	if err != nil {
		return nil, fmt.Errorf("fail to load metadata : %v", err)
	}

	log.Printf("Metadata %v\n", metadata)

	// Apply arguments to the file
	workflowArguments, err := metadata.BuildWorkflowArgument(arguments)
	if err != nil {
		return nil, fmt.Errorf("fail to compile workflow parameters : %v", err)
	}

	log.Printf("Arguments : %v", workflowArguments)

	format, err := DecodeWorkflowFile(fileContent, workflowArguments)
	if err != nil {
		return nil, fmt.Errorf("fail to decode workflow : %v", err)
	}

	return BuildWorkflow(*format, *metadata)
}

// CastParameter Cast a string to the given WorkflowParameterType
func CastParameter(stringValue string, paramType WorkflowParameterType) (interface{}, error) {
	switch paramType {
	case BooleanType:
		return strconv.ParseBool(stringValue)
	case FloatType:
		return strconv.ParseFloat(stringValue, 64)
	case NumberType:
		return strconv.ParseInt(stringValue, 10, 64)
	default:
		return stringValue, nil
	}
}

// BuildWorkflowArgument creates the attributes map that can be used to
func (metadata *WorkflowMetadataTemplate) BuildWorkflowArgument(arguments map[string]string) (map[string]interface{}, error) {
	// Apply parameters to the file
	workflowArguments := make(map[string]interface{})

	for k, v := range metadata.Parameters {
		stringValue, parameterSpecified := arguments[k]
		// Get value
		if !parameterSpecified || stringValue == "" {
			if v.DefaultValue != "" {
				stringValue = v.DefaultValue
			}
			return nil, fmt.Errorf("parameter [%s] is required", k)
		}

		// Cast to the right type
		realValue, err := CastParameter(stringValue, v.Type)
		if err != nil {
			return nil, fmt.Errorf("fail to cast parameter [%s] with value [%s] to type [%s] : %v", k, stringValue, v.Type, err)
		}

		// Add to the parameter list
		workflowArguments[k] = realValue
	}

	return workflowArguments, nil
}

// DecodeWorkflowFile Decode the workflow file using the given context and return the result
func DecodeWorkflowFile(fileContent []byte, workflowArguments map[string]interface{}) (*WorkflowFileTemplate, error) {
	// Parse the file with arguments
	decodedFile, err := ApplyContext(fileContent, workflowArguments)

	yamlDecoder := yaml.NewDecoder(bytes.NewReader(decodedFile))

	// Skip the metadata part of the file
	err = yamlDecoder.Decode(make(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("fail to parse workflow : %v", err)
	}

	var workflowData WorkflowFileTemplate

	err = yamlDecoder.Decode(&workflowData)
	if err != nil {
		return nil, fmt.Errorf("fail to parse workflow : %v", err)
	}

	log.Printf("Workflow : %v", workflowData)

	return &workflowData, nil
}

// ApplyContext Format the workflow file with the context to decode the go templates.
// It returns the decoded file
func ApplyContext(fileContent []byte, workflowArguments map[string]interface{}) ([]byte, error) {
	// Parse the file with arguments
	tmpl, err := template.New("configParser").Parse(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("fail to parse template %v", err)
	}

	buf := &bytes.Buffer{}

	context := struct {
		Parameters map[string]interface{}
	}{workflowArguments}

	err = tmpl.Execute(buf, context)
	if err != nil {
		return nil, fmt.Errorf("fail to parse workflow : %v", err)
	}

	return buf.Bytes(), nil
}
