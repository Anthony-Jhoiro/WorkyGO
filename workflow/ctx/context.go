package ctx

import "github.com/Anthony-Jhoiro/WorkyGO/logger"

type WorkflowContext interface {
	GetExternalTemplate(string) ([]byte, error)
	GetLogger() logger.Logger
	SetLogger(logger.Logger)
	Copy() WorkflowContext
	GetRunNumber() string
}
