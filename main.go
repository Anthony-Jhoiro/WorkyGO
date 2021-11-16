package main

import (
	"Workflow/config"
	"Workflow/docker"
	"log"
)

func main() {

	//conf := config.ReadConfFile("config.yaml")
	//
	//name, err := conf.Name()
	//if err != nil {
	//	log.Fatalf("Err %v", err)
	//}
	//
	//log.Fatalf("=> %s", name)

	cnf, err := config.LoadWorkflowConfig("./examples/example1")
	if err != nil {
		log.Fatalf("Config parse error %v", err)
	}

	log.Printf("=> %v", cnf)

	for _, step := range cnf.Workflow.Steps {
		dockerConfig := &docker.DockerImageConfig{
			Image:   step.Image,
			Command: step.Commands,
			Config: docker.Config{
				Entrypoint: "/bin/sh",
				Name:       "tata",
				Commands:   []string{step.Commands},
			},
		}

		container, err := docker.NewContainer(dockerConfig)
		if err != nil {
			log.Fatalf("Fail to create container %v\n", err)
		}

		err = container.Init()
		if err != nil {
			log.Fatalf("Fail to initialise container %v\n", err)
		}

		err = container.Run()
		if err != nil {
			log.Fatalf("Fail to run the container %v\n", err)
		}

		logs, err := container.GetLogs()
		if err != nil {
			log.Fatalf("Fail to parse logs %v\n", err)
		}

		log.Print(logs)
	}

	//config := &docker.DockerImageConfig{
	//	Image:   "docker/whalesay",
	//	Command: "cowsay boo",
	//	Config: docker.Config{
	//		Volumes: []docker.VolumeConfig{
	//			{
	//				Label:            "Test",
	//				ContainerMapping: "/app",
	//				ReadOnly:         true,
	//				Persistent:       false,
	//			},
	//		},
	//	},
	//}
	//
	//err := config.Run()
	//if err != nil {
	//	log.Println(err)
	//}

	//stepA := workflow.NewSimpleStep("A", "A", "Lorem Ipsum", waitAndSay)
	//stepB := workflow.NewSimpleStep("B", "B", "Lorem Ipsum", waitAndSay)
	//stepC := workflow.NewSimpleStep("C", "C", "Lorem Ipsum", waitAndSay)
	//stepD := workflow.NewSimpleStep("D", "D", "Lorem Ipsum", waitAndSay)
	//
	//stepB.AddRequirement(stepA)
	//stepC.AddRequirement(stepA)
	//stepD.AddRequirement(stepC)
	//stepD.AddRequirement(stepB)
	//
	//wf := workflow.NewWorkflow(stepA, []*workflow.Step{stepA, stepB, stepC, stepD})
	//
	//wf.Run()
}
