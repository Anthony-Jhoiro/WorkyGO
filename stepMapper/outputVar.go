package stepMapper

import "fmt"

func getOutputParser(output map[string]map[string]string) func(string, string) (string, error) {
	return func(stepName, varName string) (string, error) {
		step, ok := output[stepName]
		if !ok {
			return "", fmt.Errorf("step %s not found", stepName)
		}

		value, ok := step[varName]
		if !ok {
			return "", fmt.Errorf("output value %s not found", varName)
		}

		return value, nil
	}
}
