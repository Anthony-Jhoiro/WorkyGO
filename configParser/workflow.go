package configParser

import (
	"Workflow/workflow"
)

type Runner struct {
	Name        string
	Description string
	Maintainer  string
	Arguments   string
	Workflow    *workflow.Workflow
}
