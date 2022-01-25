package configParser

//func ParseWorkflowTemplate(ctx context.Context, fileContent []byte, arguments map[string]string) (WorkflowContext, error) {
//	// Load the file metadata
//	metadata, err := LoadMetadata(fileContent)
//	if err != nil {
//		return nil, fmt.Errorf("fail to load metadata : %v", err)
//	}
//
//	if len(metadata.Imports) > 0 {
//		log.Println("[WARNING] : Using external templates can expose your system to several risks.")
//		for _, externalTemplate := range metadata.Imports {
//			err = ResolveExternalTemplate(externalTemplate)
//			if err != nil {
//				return nil, fmt.Errorf("fail to import template : %v", err)
//			}
//		}
//	}
//
//
//	logger.LOG.Debug(fmt.Sprintf("Metadata %v", metadata))
//
//	// Apply arguments to the file
//	workflowArguments, err := metadata.BuildWorkflowArgument(arguments)
//	if err != nil {
//		return nil, fmt.Errorf("fail to compile workflow parameters : %v", err)
//	}
//
//	logger.LOG.Debug(fmt.Sprintf("Arguments : %v", workflowArguments))
//
//	format, err := DecodeWorkflowFile(fileContent, workflowArguments)
//	if err != nil {
//		return nil, fmt.Errorf("fail to decode workflow : %v", err)
//	}
//
//	logger.LOG.Debug(fmt.Sprintf("Workflow : %v", format))
//
//	return BuildWorkflowContext(ctx, metadata, nil, workflowArguments), nil
//}
