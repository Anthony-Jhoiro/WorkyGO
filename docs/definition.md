# Definition

The definition part of the workflow file contains the step definitions.

Example :
```yaml
workflow:
  steps:
    - name: step-1
      image: alpine:latest
      commands: echo Step 1

    {{ if (index .Parameters "run_step_2") }}

    - name: step-2
      image: node:lts
      commands: echo Step 2
    {{ end }}
```

The available steps are :
- [Docker Step](/docs/docker-step.md)
- [Workflow Step](/docs/workflow-step.md)
