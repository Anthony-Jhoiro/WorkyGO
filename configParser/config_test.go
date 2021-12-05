package configParser

import "testing"

var configExample = []byte(`
name: My magnificent workflow
parameters:
  - name: arg1
    description: first argument
    validators: 'valid1'
    default_value: toto
    type: string
`)

func TestParseConfigFileContent(t *testing.T) {
	template, err := ParseConfigFileContent(configExample)

	if err != nil {
		t.Errorf("Parsing Error %v", err)
	}

	if template == nil ||
		template.Name != "My magnificent workflow" ||
		len(template.Parameters) != 1 ||
		template.Parameters[0].Name != "arg1" ||
		template.Parameters[0].DefaultValue != "toto" ||
		template.Parameters[0].Type != StringType ||
		template.Parameters[0].Description != "first argument" ||
		template.Parameters[0].Validators != "valid1" {

		t.Errorf("Parsing Error : Invalid template")
	}

}
