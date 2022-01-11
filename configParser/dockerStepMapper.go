package configParser

import (
	"Workflow/docker"
	"Workflow/workflow"
	"fmt"
	"log"
)

func (tpl *StepDockerFormat) ToWorkFlowStep() *workflow.Step {
	return workflow.NewSimpleStep(tpl.Name, tpl.Name, "", tpl.Run)
}

func (tpl *StepDockerFormat) Run(context *interface{}) error {

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
	// Pull image

	container, err := docker.NewContainer(dockerConfig)
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

	logs, err := container.GetLogs()
	if err != nil {
		log.Fatalf("Fail to parse logs %v\n", err)
	}
	log.Println(logs)

	return nil
}
