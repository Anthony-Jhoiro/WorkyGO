package docker

import (
	"context"
	"fmt"
	"github.com/Anthony-Jhoiro/WorkyGO/workflow/ctx"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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
func (c *Container) Init(ctx ctx.WorkflowContext) error {
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
	err = c.PullImage(ctx)
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

func (c *Container) PullImage(ctx ctx.WorkflowContext) error {
	// Pull docker image
	out, err := c.client.ImagePull(c.Context, c.config.image(), types.ImagePullOptions{})

	if err != nil {
		panic(err)
	}

	getLogger := ctx.GetLogger()
	_, err = getLogger.PrintFormattedReader(0, "PULL - %s", out)
	if err != nil {
		return fmt.Errorf("fail to redirect logs")
	}

	err = out.Close()
	return err
}

func (c *Container) Run(ctx ctx.WorkflowContext) (map[string]string, error) {

	output, execErr := c.exec(ctx)

	// Clean the process
	clearErr := c.clear()
	if clearErr != nil {
		return nil, clearErr
	}

	return output, execErr
}

// exec execute the main process of the container
func (c *Container) exec(ctx ctx.WorkflowContext) (map[string]string, error) {
	// Start the container
	err := c.client.ContainerStart(c.Context, c.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
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
		return nil, err
	}

	getLogger := ctx.GetLogger()
	output, err := getLogger.PrintFormattedReader(8, "%s", out)
	if err != nil {
		return nil, fmt.Errorf("fail to redirect logs : %v", err)
	}

	// Wait until the container exec process is completed
	statusCh, errCh := c.client.ContainerWait(c.Context, c.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case a := <-statusCh:
		if a.StatusCode != 0 {
			return nil, fmt.Errorf("step failed with status %d", a.StatusCode)
		}
	}
	return output, nil

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
