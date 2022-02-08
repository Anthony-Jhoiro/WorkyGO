package configParser

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"strconv"
	"text/template"
)

// ParseWorkflowFile Parse the content of a workflow file with the given arguments
func ParseWorkflowFile(fileContent []byte, arguments map[string]string) (*ParsedWorkflow, error) {
	// Load the file metadata
	metadata, err := loadMetadata(fileContent)
	if err != nil {
		return nil, fmt.Errorf("fail to load metadata : %v", err)
	}

	externalTemplates, err := ResolveExternalTemplates(metadata.Imports)
	if err != nil {
		return nil, fmt.Errorf("fail during external templates importation : %v", err)
	}

	// Apply arguments to the file
	workflowArguments, err := metadata.buildWorkflowArgument(arguments)
	if err != nil {
		return nil, fmt.Errorf("fail to compile workflow parameters : %v", err)
	}

	format, err := decodeWorkflowFile(fileContent, workflowArguments)
	if err != nil {
		return nil, fmt.Errorf("fail to decode workflow : %v", err)
	}

	return &ParsedWorkflow{
		Name:        metadata.Name,
		Description: metadata.Description,
		Maintainer:  metadata.Maintainer,
		Steps:       format.Workflow.Steps,
		Imports:     externalTemplates,
		Output:      metadata.Output,
	}, nil
}

// castParameter Cast a string to the given WorkflowParameterType
func castParameter(stringValue string, paramType WorkflowParameterType) (interface{}, error) {
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

// buildWorkflowArgument creates the attributes map that can be used to
func (metadata *workflowMetadataTemplate) buildWorkflowArgument(arguments map[string]string) (map[string]interface{}, error) {
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
		realValue, err := castParameter(stringValue, v.Type)
		if err != nil {
			return nil, fmt.Errorf("fail to cast parameter [%s] with value [%s] to type [%s] : %v", k, stringValue, v.Type, err)
		}

		// Add to the parameter list
		workflowArguments[k] = realValue
	}

	return workflowArguments, nil
}

// decodeWorkflowFile Decode the workflow file using the given context and return the result
func decodeWorkflowFile(fileContent []byte, workflowArguments map[string]interface{}) (*workflowFileTemplate, error) {
	// Parse the file with arguments
	decodedFile, err := applyContext(fileContent, workflowArguments)

	yamlDecoder := yaml.NewDecoder(bytes.NewReader(decodedFile))

	// Skip the metadata part of the file
	err = yamlDecoder.Decode(make(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("fail to parse workflow : %v", err)
	}

	var workflowData workflowFileTemplate

	err = yamlDecoder.Decode(&workflowData)
	if err != nil {
		return nil, fmt.Errorf("fail to parse workflow : %v", err)
	}

	return &workflowData, nil
}

func ReferenceStepOutput(stepName, varName string) string {
	return fmt.Sprintf("{{ getVar \"%s\" \"%s\" }}", stepName, varName)
}

// applyContext Format the workflow file with the context to decode the go templates.
// It returns the decoded file
func applyContext(fileContent []byte, workflowArguments map[string]interface{}) ([]byte, error) {
	// Parse the file with arguments
	tmpl, err := template.New("configParser").Funcs(template.FuncMap{"getVar": ReferenceStepOutput}).Parse(string(fileContent))
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
