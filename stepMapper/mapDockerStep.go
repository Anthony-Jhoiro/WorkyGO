package stepMapper

import (
	"Workflow/docker"
	"Workflow/workflow/ctx"
	"encoding/json"
	"fmt"
	"strings"
)

type StepDockerPersistFormat struct {
	Name   string
	Source string
}

type StepDocker struct {
	Id        string                    `json:"id,omitempty"`
	Name      string                    `json:"name"`
	Image     string                    `json:"image"`
	Workdir   string                    `json:"workdir,omitempty"`
	Commands  string                    `json:"commands"`
	DependsOn []string                  `json:"depends_on,omitempty"`
	Persist   []StepDockerPersistFormat `json:"persist,omitempty"`
	Volumes   []docker.VolumeConfig     `json:"volumes,omitempty"`
}

func (ds *StepDocker) GetDependencies() []string {
	if ds.DependsOn == nil {
		return make([]string, 0)
	}
	return ds.DependsOn
}

func MapDockerStep(template interface{}) (*StepDocker, error) {

	jsonBytes, err := json.Marshal(template)
	if err != nil {
		return nil, fmt.Errorf("fail to parse step : %v", err)
	}

	var dockerStep StepDocker

	err = json.Unmarshal(jsonBytes, &dockerStep)
	if err != nil {
		return nil, fmt.Errorf("invalid docker step")
	}

	if len(dockerStep.Commands) == 0 {
		return nil, fmt.Errorf("the attribute 'commands' requires a value")
	}

	return &dockerStep, nil
}

func (ds *StepDocker) Init(_ ctx.WorkflowContext) error {
	return nil
}

func (ds *StepDocker) Clean() {

}

func (ds *StepDocker) GetLabel() string {
	return strings.ToLower(strings.ReplaceAll(ds.Name, " ", "_"))
}

func (ds *StepDocker) GetName() string {
	return ds.Name
}

func (ds *StepDocker) GetDescription() string {
	return ""
}

func (ds *StepDocker) Run(ctx ctx.WorkflowContext) error {

	volumes := make([]docker.VolumeConfig, 0, len(ds.Persist))

	for _, v := range ds.Persist {
		volumes = append(volumes, docker.VolumeConfig{
			Label:            v.Name,
			ContainerMapping: v.Source,
			ReadOnly:         false,
			Persistent:       true,
		})
	}

	dockerConfig := &docker.DockerImageConfig{
		Image:   ds.Image,
		Command: ds.Commands,
		Config: docker.Config{
			Volumes:    volumes,
			Env:        nil,
			WorkingDir: "",
			Entrypoint: "/bin/sh",
			Name:       "tata",
			Commands:   []string{ds.Commands},
		},
	}

	container, err := docker.NewContainer(dockerConfig, ds.Id)
	if err != nil {
		return fmt.Errorf("fail to create container %v", err)
	}

	err = container.Init(ctx)
	if err != nil {
		return fmt.Errorf("fail to initialise container %v", err)
	}

	err = container.Run(ctx)
	if err != nil {
		return fmt.Errorf("fail to run the container %v", err)
	}

	return nil
}
