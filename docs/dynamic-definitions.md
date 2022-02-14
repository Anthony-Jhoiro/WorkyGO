# Dynamic definitions

Your workflow definitions cn be dynamic depending on the parameters that you pass to your workflow. You can access parameters values with the map `.Parameters` that contains the values of all parsed parameters. 

To see a list of all available possibilities, see the [documentation of go templates](https://pkg.go.dev/text/template#hdr-Actions)



```yaml
name: Dynamic Workflow Example
description: lorem ipsum dolor sit ahmet

maintainer: Anthony-Jhoiro

parameters:
  run_step_2:
    type: boolean
---

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

