package configParser

import (
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
