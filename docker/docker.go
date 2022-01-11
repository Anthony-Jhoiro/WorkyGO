package docker

import (
	"github.com/docker/docker/api/types/container"
)

func createContainer(c Container) (container.ContainerCreateCreatedBody, error) {

	return c.client.ContainerCreate(c.Context, &container.Config{
		AttachStdin:  false,
		AttachStdout: false,
		AttachStderr: false,
		Env:          c.config.getEnvVars(),
		Image:        c.config.image(),
		WorkingDir:   c.config.getWorkingDirectory(),
		Entrypoint:   []string{"/bin/sh", "-c", "/entrypoint/entrypoint.sh"},
		StopSignal:   "", // Work with that
		StopTimeout:  c.config.getTimeout(),
	}, &container.HostConfig{
		Mounts: c.config.getDockerVolumes(),
	}, nil, nil, c.Name)
}
