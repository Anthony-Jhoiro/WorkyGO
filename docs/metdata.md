# Metadata

## `name` (string)
Name of your workflow to find it easily
```yaml
name: My Super Awesome Workflow
```
## `description` (string)
Short description of your workflow to allow anyone to understand its usages and behavior
```yaml
description: Do this and that and return that
---
description: |
  This is a 
  multiline 
  description
```
## `maintainer` (string)
Name of the maintainer of the workflow
```yaml
maintainer: Alan Turing (alan.turing@mail.com)
```

## `parameters` (map)
Map that describes the arguments needed to run your workflow. These arguments can be used when the go template is parsed in the definition.

```yaml
parameters:
  foo:
    description: Argument example
    type: string
    
  bar:
    description: Argument example with default value
    type: number

  toto:
    description: Boolean Argument
    type: boolean
```

### `parameters`/`description` (string)
Description of your argument

### `parameters`/`type` (string)
Type of the parameter, the supported types are :
- number : any integer
- float : any number with are without floating point
- string : any string
- boolean : 
  - **for true :** 1, t, T, TRUE, true, True
  - **for false :** 0, f, F, FALSE, false, False


## `imports` (list)
List of external templates used in the definition part. They can be used to split your workflow into multiple smaller workflows.

> WARNING : The usage of template can expose your system to several risks, use them with precaution.

```yaml
imports: 
  - name: gitclone
    url: https://raw.githubusercontent.com/Anthony-Jhoiro/sample_git/master/test_action.yml
```

### `imports`/`name` (string)
Name of the workflow, the name will be used to reference the workflow later in the definition.

### `imports`/`url` (string)
Url of the file to import.

> Tip : If the template that you want to import is located in a GitHub repository, you can use `githubusercontent.com` to access it.