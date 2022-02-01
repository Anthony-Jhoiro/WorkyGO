package configParser

import (
	"Workflow/logger"
	"Workflow/workflow/ctx"
	"fmt"
	"io/ioutil"
)

type ParsedWorkflow struct {
	Name        string
	Description string
	Maintainer  string
	Arguments   string
	Steps       []interface{}
	Imports     map[string]string
	log         logger.InteractiveLogger
}

func (wf *ParsedWorkflow) GetExternalTemplate(tplName string) ([]byte, error) {
	fileName, ok := wf.Imports[tplName]
	if !ok {
		return nil, fmt.Errorf("template [%s] not found", tplName)
	}

	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("fail to read file %s : %v", fileName, err)
	}

	return fileContent, nil
}

func (wf *ParsedWorkflow) GetLogger() logger.InteractiveLogger {
	return wf.log
}

func (wf *ParsedWorkflow) SetLogger(log logger.InteractiveLogger) {
	wf.log = log
}

func (wf *ParsedWorkflow) Copy() ctx.WorkflowContext {
	return &ParsedWorkflow{
		Name:        wf.Name,
		Description: wf.Description,
		Maintainer:  wf.Maintainer,
		Arguments:   wf.Arguments,
		Steps:       wf.Steps,
		Imports:     wf.Imports,
		log:         wf.log,
	}
}
