package configParser

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
)

// LoadMetadata Extract the metadata part of the workflow file
//and check the required field
func LoadMetadata(fileContent []byte) (*WorkflowMetadataTemplate, error) {
	metadata := WorkflowMetadataTemplate{}

	fileReader := bytes.NewReader(fileContent)
	yamlDecoder := yaml.NewDecoder(fileReader)

	// Decode metadata
	err := yamlDecoder.Decode(&metadata)
	if err != nil {
		return nil, fmt.Errorf("fail to parse workflow metadata %v", err)
	}

	// Check Required fields
	if metadata.Name == "" {
		return nil, fmt.Errorf("your workflow must have a name attribute")
	}

	return &metadata, nil
}
