package config

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"text/template"
)

type WorkflowConfig struct {
	Data interface{}
	Env  map[string]string
}

type envFile map[string]interface{}

type StepFormat struct {
	Name     string
	Image    string
	Workdir  string
	Commands string
}

type WorkflowWorkflowFormat struct {
	Steps []StepFormat
}

type WorkflowFormat struct {
	Workflow WorkflowWorkflowFormat
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func LoadWorkflowConfig(basedir string) (*WorkflowFormat, error) {
	// Load env file

	environmentConfig := &envFile{}
	workflowConfig := &WorkflowFormat{}

	envFilePath := path.Join(basedir, "env.yaml")
	if fileExists(envFilePath) {
		envConfigContent, err := ioutil.ReadFile(envFilePath)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(envConfigContent, &environmentConfig)
		if err != nil {
			return nil, err
		}
	}

	// Load config

	workflowFilePath := path.Join(basedir, "workflow.yaml")

	if fileExists(workflowFilePath) {
		workflowConfigContent, err := ioutil.ReadFile(workflowFilePath)
		if err != nil {
			return nil, err
		}
		tmpl, err := template.New("config").Parse(string(workflowConfigContent))
		if err != nil {
			return nil, err
		}

		buf := &bytes.Buffer{}

		err = tmpl.Execute(buf, environmentConfig)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(buf.Bytes(), &workflowConfig)
		if err != nil {
			return nil, err
		}
	}

	return workflowConfig, nil
}
