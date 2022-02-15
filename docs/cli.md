# CLI

For a list of available usages, run `workygo --help`.

## Commands
For the moment, only the `run` command is supported with the basics `help` and `completion`.

The run command take as a positional argument the name of the workflow file.

### Pass arguments to workflow
To pass arguments to workflow, you can either use `-a` or `--arg` with the argument and its value in the command line like this : `-a key=value`.

You can also import arguments from a json file with `-f` or `--from-json` like this : `-f my-args.json`.

## Examples
```bash
workygo run ci.yaml -f ci-args.json
```

