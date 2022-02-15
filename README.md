# WorkyGO

**Repo** : [https://github.com/Anthony-Jhoiro/WorkyGO](https://github.com/Anthony-Jhoiro/WorkyGO)

WorkyGo is a tool to manage workflows. These workflows can be used to automate your daily tasks or build your CI/CD pipeline. Anything is possible !

You can find some working examples in the example directory. If you want to see a complete workflow example, I recommend starting with the [ci example](examples/ci.yaml).

## Authors 
- Anthony Quéré (Anthony-Jhoiro)

## Documentation
- [Use the CLI](docs/cli.md)
- [Basics](docs/usages.md)
- [Metadata](docs/metdata.md)
- [Definition](docs/definition.md)
- [Docker Step](docs/docker-step.md)
- [Workflow Step](docs/workflow-step.md)
- [Work With outputs](docs/work-with-ouputs.md)

## Functionalities
- Create workflow with parameters
- Use any linux based image to execute steps
- Use go template to dynamically change the workflow definition
- Use step output in other steps to pass variables
- Share files and folders between steps through volumes
- Execute external workflows as tests

## Next steps 
- Allow default values for parameters
- Create a registry to manage workflows
- Add environment variables to docker steps
- Create an API interface
- Secret management
- Control workflow validity with a simplex algorithm
- Support custom docker registry
- Support registry authentication
- Support for private external workflows