package ctx

type WorkflowContext interface {
	GetExternalTemplate(string) ([]byte, error)
}
