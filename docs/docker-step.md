# Docker step
The Docker step lets you create steps that will be executed inside Docker containers. The only restriction is that these containers are based on Linux and hosted on Docker Hub are that the images are already pulled and available in your local Docker images.

## Syntax

### `name` (string)
Name of the step, you can use any letters and dashes.

### `image` (string)
Name of the Docker image to use with its version (ex: alpine/git:v2.32.0)

### `workdir` (string) - optional
Base workdir for your commands execution

### `commands` (string)
Commands that will run on that container, to use multiline command, you can use the syntax with the `|` character :

```yaml
commands: |
  apt-get update -y
  apt-get install cowsay
  cowsay "Hello World"
```

### `persist` (list) - optional
List of volume that will be persisted through the steps

#### `persist`/`name` (string)
Name of the volume (it won't be called that way in when created in Docker, but you will be able to reference it  with that name)

#### `persist`/`source` (string)
Path in the container that will be mapped to that volume.

### `depends_on` (list)
List of `string` that represents the steps that needs to be executed before that one.

