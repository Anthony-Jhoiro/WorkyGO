package ctx

import "Workflow/logger"

type WorkflowContext interface {
	GetExternalTemplate(string) ([]byte, error)
	GetLogger() logger.Logger
	SetLogger(logger.Logger)
	Copy() WorkflowContext
}
