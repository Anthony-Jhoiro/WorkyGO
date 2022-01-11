package configParser

import (
	"Workflow/docker"
	"Workflow/workflow"
	"fmt"
)

func (tpl *StepDocker) ToWorkFlowStep() *workflow.Step {
	return workflow.NewSimpleStep(tpl.Name, tpl.Name, "", tpl.Run)
}

func (tpl *StepDocker) Run(_ *interface{}) error {

	volumes := make([]docker.VolumeConfig, 0, len(tpl.Persist))

	for _, v := range tpl.Persist {
		volumes = append(volumes, docker.VolumeConfig{
			Label:            v.Name,
			ContainerMapping: v.Source,
			ReadOnly:         false,
			Persistent:       true,
		})
	}

	dockerConfig := &docker.DockerImageConfig{
		Image:   tpl.Image,
		Command: tpl.Commands,
		Config: docker.Config{
			Volumes:    volumes,
			Env:        nil,
			WorkingDir: "",
			Entrypoint: "/bin/sh",
			Name:       "tata",
			Commands:   []string{tpl.Commands},
		},
	}

	container, err := docker.NewContainer(dockerConfig, tpl.Id)
	if err != nil {
		return fmt.Errorf("fail to create container %v", err)
	}

	err = container.Init()
	if err != nil {
		return fmt.Errorf("fail to initialise container %v", err)
	}

	err = container.Run()
	if err != nil {
		return fmt.Errorf("fail to run the container %v", err)
	}

	return nil
}
