package ctx

import "Workflow/logger"

type WorkflowContext interface {
	GetExternalTemplate(string) ([]byte, error)
	GetLogger() logger.InteractiveLogger
	SetLogger(logger.InteractiveLogger)
	Copy() WorkflowContext
}
