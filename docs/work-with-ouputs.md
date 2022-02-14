# Work with outputs

Steps can emit output values that can be used as variables during your workflow execution. 

## Return Output from Docker Step
To add output to a Docker Step, you need to log it in the standard output (stdout). The log must be on one line and be formatted like this : `^::output::key::value::$`

**Example :**
```
log 1
log 2
::output::foo::bar::
log 3
```

## Return Output from workflow
In your workflow metadata you need to add :
```yaml
output:
  <keyName>: <stepName><keyValue>
```

**Example :**
```yaml
output:
  foo: step-1.bar
```

> It also how you emit output using workflow steps

## Use output in Steps
To use output values in steps, you need to use `{{ getVar "step-name" "key-name" }}` in the values. But **you can't do that in all fields !**

In docker steps, you can use them for the fields `commands` and `workdir`

In workflow steps, you can use them for `parameters` values

```yaml
commands: echo {{ getVar "step-1" "foo" }}
```
