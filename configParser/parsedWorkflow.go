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
	Steps       []interface{}
	Imports     map[string]string
	Output      map[string]string
	log         logger.Logger
	runNumber   string
}

func (wf *ParsedWorkflow) GetRunNumber() string {
	return wf.runNumber
}

func (wf *ParsedWorkflow) SetRunNumber(rn string) {
	wf.runNumber = rn
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

func (wf *ParsedWorkflow) GetLogger() logger.Logger {
	return wf.log
}

func (wf *ParsedWorkflow) SetLogger(log logger.Logger) {
	wf.log = log
}

func (wf *ParsedWorkflow) Copy() ctx.WorkflowContext {
	return &ParsedWorkflow{
		Name:        wf.Name,
		Description: wf.Description,
		Maintainer:  wf.Maintainer,
		Steps:       wf.Steps,
		Imports:     wf.Imports,
		log:         wf.log,
		runNumber:   wf.runNumber,
	}
}
