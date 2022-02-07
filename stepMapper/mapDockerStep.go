package stepMapper

import (
	"Workflow/docker"
	"Workflow/workflow/ctx"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	"log"
	"strings"
	"text/template"
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
	Executor  string                    `json:"executor,omitempty"`
	volumes   []docker.VolumeConfig
	client    *client.Client
	output    map[string]string
}

func (ds *StepDocker) GetOutput() map[string]string {
	return ds.output
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

func _resolveStepValue(parser *template.Template, value *string) error {

	tpl, err := parser.Parse(*value)
	if err != nil {
		return fmt.Errorf("parsing error : %v", err)
	}

	buf := &bytes.Buffer{}

	err = tpl.Execute(buf, nil)
	if err != nil {
		return fmt.Errorf("parsing error : %v", err)
	}

	*value = buf.String()

	return nil
}

func (ds *StepDocker) resolveStepValues(previousStepsOutput map[string]map[string]string) error {
	parser := template.New("parser").Funcs(template.FuncMap{"getVar": getOutputParser(previousStepsOutput)})

	err := _resolveStepValue(parser, &ds.Workdir)
	if err != nil {
		return fmt.Errorf("fail to parse workdir : %v", err)
	}
	err = _resolveStepValue(parser, &ds.Commands)
	if err != nil {
		return fmt.Errorf("fail to parse commands : %v", err)
	}

	return nil
}

func (ds *StepDocker) Init(ctx ctx.WorkflowContext, previousStepsOutput map[string]map[string]string) error {

	err := ds.resolveStepValues(previousStepsOutput)
	if err != nil {
		return fmt.Errorf("fail to parse some values : %v", err)
	}

	volumes := make([]docker.VolumeConfig, 0, len(ds.Persist))

	for _, v := range ds.Persist {
		volumes = append(volumes, docker.VolumeConfig{
			Label:            fmt.Sprintf("wf-%s-%s", v.Name, ctx.GetRunNumber()),
			ContainerMapping: v.Source,
			ReadOnly:         false,
			Persistent:       true,
		})
	}
	ds.volumes = volumes

	// Add docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("fatal : fail to create docker client : %v", err)
	}

	ds.client = cli

	return nil
}

func (ds *StepDocker) Clean() {

	c := context.Background()

	for _, volume := range ds.volumes {
		inspect, err := ds.client.VolumeInspect(c, volume.Label)
		if err == nil {
			err := ds.client.VolumeRemove(c, inspect.Name, true)
			if err != nil {
				log.Printf("[WARNING] Fail to delete volume [%s] : %v", inspect.Name, err)
			}
		}
	}
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

	executor := "/bin/sh -c"
	if len(ds.Executor) > 0 {
		executor = ds.Executor
	}

	dockerConfig := &docker.DockerImageConfig{
		Image:   ds.Image,
		Command: ds.Commands,
		Config: docker.Config{
			Volumes:    ds.volumes,
			Env:        nil,
			WorkingDir: "",
			Entrypoint: executor,
			Name:       ds.Name,
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

	output, err := container.Run(ctx)
	if err != nil {
		return fmt.Errorf("fail to run the container %v", err)
	}

	ds.output = output

	return nil
}
