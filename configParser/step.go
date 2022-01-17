package configParser

import (
	"encoding/json"
	"fmt"
)

func MapDictToStep(stepDefinition interface{}) (WorkflowStep, error) {
	jsonBytes, err := json.Marshal(stepDefinition)
	if err != nil {
		return nil, fmt.Errorf("fail to parse step : %v", err)
	}

	var dockerStep StepDockerFormat

	err = json.Unmarshal(jsonBytes, &dockerStep)
	if err == nil {
		return dockerStep, nil
	}

	importedStep, ok := stepDefinition.(StepImportedFormat)
	if ok {
		return importedStep, nil
	}

	return nil, fmt.Errorf("invalid step")
}
