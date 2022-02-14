# How to use

## With CLI

### Write a workflow file

A workflow file is a Yaml file composed of a **Metadata** part and a **Definition** Part. You can use [go templates](https://pkg.go.dev/text/template#hdr-Actions) for the second part. The metadata part contains all infos needed to start and reference your workflow and the definition part contains all elements to run your workflow.

### Definition

#### `workflow`/`steps` (list)
List of the different steps tht compose the workflow. At the moment there are 2 type of steps :
- [Docker step](docker-step.md)

