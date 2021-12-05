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
	dockerConfig := &docker.DockerImageConfig{
		Image:   tpl.Image,
		Command: tpl.Commands,
		Config: docker.Config{
			Entrypoint: "/bin/sh",
			Name:       "tata",
			Commands:   []string{tpl.Commands},
		},
	}

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
