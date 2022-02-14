package stepMapper

import (
	"bytes"
	"fmt"
	"text/template"
)

func _resolveStepValue(parser *template.Template, value *string) error {

	tpl, err := parser.Parse(*value)
	if err != nil {
		return fmt.Errorf("parsing error : %v", err)
	}

	buf := &bytes.Buffer{}

	err = tpl.Execute(buf, nil)
	if err != nil {
		return fmt.Errorf("parsing error : %v", err)
	}

	*value = buf.String()

	return nil
}

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
