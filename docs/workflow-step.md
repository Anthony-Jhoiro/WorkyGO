# Workflow step

The workflow step allows you to execute an external workflow as a Step in your own workflow. For example, lets say that you want to reproduce a part of your workflow into multiple other workflow. You can buld a separate workflow with only the shared part, add it into a GitHub repository for example and import it from all of your workflows. This way, the logic is located in a single place and you can focus on the specific part of your workflow.

External workflows work with `imports`, see the declaration in the metadata syntax

## Syntax

### `name` (string)
Name of the step, you can use any letters and dashes.

### `workflow` (string)
Name of the imported workflow (see the declarations of `imports` in the metadata syntax)

### `parameters` (map\[string\]string) - optional
Arguments passed to the parameters of the external workflow.

