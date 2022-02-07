# How to use

## With CLI

### Write a workflow file

A workflow file is a Yaml file composed of a **Metadata** part and a **Definition** Part. You can use [go templates](https://pkg.go.dev/text/template#hdr-Actions) for the second part. The metadata part contains all infos needed to start and reference your workflow and the definition part contains all elements to run your workflow.

#### Metadata

##### `name` (string)
Name of your workflow to find it easily
##### `description` (string)
Short description of your workflow to allow anyone to understand its usages and behavior
##### `maintainer` (string)
Name of the maintainer of the workflow

##### `parameters` (map)
Map that describes the arguments needed to run your workflow. These arguments can be used when the go template is parsed in the definition.

###### `parameters`/`description` (string)
Description of your argument

###### `parameters`/`default_value` (string)
Default value that will be passed to the parameter if no value is specified.

###### `parameters`/`type` (string)
Type of the parameter, the supported types are :
- number : any integer
- float : any number with are without floating point
- string : any string
- boolean : if `yes` or `true` is passed, it will be true, otherwise it will be false


##### `imports` (list)
List of external templates used in the definition part. They can be used to split your workflow into multiple smaller workflows. 

> WARNING : The usage of template can expose your system to several risks, use them with precaution.

##### `imports`/`name` (string)
Name of the workflow, the name will be used to reference the workflow later in the definition.

##### `imports`/`url` (string)
Url of the file to import. 

> Tip : If the template that you want to import is located in a GitHub repository, you can use `githubusercontent.com` to access it.


### Definition
