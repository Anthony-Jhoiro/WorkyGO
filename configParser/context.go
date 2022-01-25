package configParser

//
//import (
//	"Workflow/workflow"
//	"context"
//	"time"
//)
//
//type WorkflowContext interface {
//	context.Context
//	GetMetadata() map[string]interface{}
//	GetSteps() []workflow.Step
//	GetArguments() map[string]interface{}
//}
//
//func BuildWorkflowContext(ctx context.Context, metadata map[string]interface{}, steps map[string]workflow.Step, arguments map[string]interface{}) WorkflowContext {
//	return workflowContextImpl{
//		ctx:       context.Background(),
//		metadata:  nil,
//		steps:     nil,
//		arguments: nil,
//	}
//}
//
//type workflowContextImpl struct {
//	ctx       context.Context
//	metadata  map[string]interface{}
//	steps     []workflow.Step
//	arguments map[string]interface{}
//}
//
//func (ctx workflowContextImpl) Deadline() (deadline time.Time, ok bool) {
//	return ctx.Deadline()
//}
//
//func (ctx workflowContextImpl) Done() <-chan struct{} {
//	return ctx.Done()
//}
//
//func (ctx workflowContextImpl) Err() error {
//	return ctx.Err()
//}
//
//func (ctx workflowContextImpl) Value(key interface{}) interface{} {
//	return ctx.Value(key)
//}
//
//func (ctx workflowContextImpl) GetMetadata() map[string]interface{} {
//	return ctx.metadata
//}
//func (ctx workflowContextImpl) GetSteps() []workflow.Step {
//	return ctx.steps
//}
//func (ctx workflowContextImpl) GetArguments() map[string]interface{} {
//	return ctx.arguments
//}
