package docker

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"os"
)

type VolumeConfig struct {
	Label            string
	ContainerMapping string
	ReadOnly         bool
	Persistent       bool
}

type Config struct {
	Volumes        []VolumeConfig
	Env            map[string]string
	WorkingDir     string
	Entrypoint     string
	Name           string
	Commands       []string
	entrypointFile *os.File
}

type DockerImageConfig struct {
	Image   string
	Command string
	Config
}

type Container struct {
	StepId string
	config configInterface
	container.ContainerCreateCreatedBody
	client *client.Client
	context.Context
	Name string
}

type configInterface interface {
	getEnvVars() []string            // returns the list of environment variables usable byt the container
	image() string                   // return the docker image to use to create the container
	getWorkingDirectory() string     // return the working directory used during the container execution
	getTimeout() *int                // return the maximum time that a container can live
	getDockerVolumes() []mount.Mount // return the list of volumes that needs to be mounted into the container
	initialize(c Container) error    // initialize function of the configuration, run during the initialisation of the container, before its creation
	destroy() error                  // destroy function of the configuration that will be called during the container destruction
	getEntrypoint() []string
}
