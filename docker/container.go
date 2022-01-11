package docker

import (
	"Workflow/logger"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"path"
	"time"
)

// RandContainerName generate a random name for a docker container base on a timestamp
func RandContainerName() string {
	now := time.Now()
	return fmt.Sprintf("c-%d", now.UnixNano())
}

// NewContainer create a new container instance from a configuration
func NewContainer(config *DockerImageConfig, stepId string) (*Container, error) {
	c := &Container{}
	c.StepId = stepId
	c.config = config

	return c, nil
}

// Init initialize the container and its configuration
func (c *Container) Init() error {
	// Setup Context
	c.Context = context.Background()

	var err error
	// Add docker client
	c.client, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// Container name
	c.Name = RandContainerName()

	// Pull container image
	err = c.PullImage()
	if err != nil {
		return err
	}

	// Set up the configuration
	err = c.config.initialize(*c)
	if err != nil {
		return err
	}

	// Create the container
	c.ContainerCreateCreatedBody, err = createContainer(*c)

	return err
}

func (c *Container) PullImage() error {
	// Pull docker image
	out, err := c.client.ImagePull(c.Context, c.config.image(), types.ImagePullOptions{})

	if err != nil {
		panic(err)
	}

	err = logger.LOG.Log(path.Join("steps", c.StepId, "pull.log"), out)
	if err != nil {
		return fmt.Errorf("fail to redirect logs")
	}

	err = out.Close()
	return err
}

func (c *Container) Run() error {

	err := c.exec()
	if err != nil {
		return err
	}

	// Clean the process
	err = c.clear()
	if err != nil {
		return err
	}

	return nil
}

// exec execute the main process of the container
func (c *Container) exec() error {
	// Start the container
	err := c.client.ContainerStart(c.Context, c.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	// Redirect the container logs to the logger file
	out, err := c.client.ContainerLogs(c.Context, c.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      "",
		Until:      "",
		Timestamps: false,
		Follow:     true,
		Tail:       "",
		Details:    false,
	})
	if err != nil {
		return err
	}

	err = logger.LOG.Log(path.Join("steps", c.StepId, "run.log"), out)
	if err != nil {
		return fmt.Errorf("fail to redirect logs : %v", err)
	}

	// Wait until the container exec process is completed
	statusCh, errCh := c.client.ContainerWait(c.Context, c.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}
	return nil

}

// clear clean the container runtime elements :
func (c *Container) clear() error {

	// Delete entrypoint file
	err := c.config.destroy()
	if err != nil {
		return err
	}

	// Delete the container
	return c.client.ContainerRemove(c.Context, c.ID, types.ContainerRemoveOptions{
		RemoveVolumes: false,
		RemoveLinks:   false,
		Force:         false,
	})
}
