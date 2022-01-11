package main

import (
	"Workflow/configParser"
	"io/ioutil"
	"log"
)

func main() {

	//cnf, err := configParser.LoadWorkflowConfig("./examples/example1")
	//if err != nil {
	//	log.Fatalf("Config parse error %v", err)
	//}
	//
	//log.Printf("=> %v", cnf)
	//
	//for _, step := range cnf.Workflow.Steps {
	//	dockerConfig := &docker.DockerImageConfig{
	//		Image:   step.Image,
	//		Command: step.Commands,
	//		Config: docker.Config{
	//			Entrypoint: "/bin/sh",
	//			Name:       "tata",
	//			Commands:   []string{step.Commands},
	//		},
	//	}
	//
	//	container, err := docker.NewContainer(dockerConfig)
	//	if err != nil {
	//		log.Fatalf("Fail to create container %v\n", err)
	//	}
	//
	//	err = container.Init()
	//	if err != nil {
	//		log.Fatalf("Fail to initialise container %v\n", err)
	//	}
	//
	//	err = container.Run()
	//	if err != nil {
	//		log.Fatalf("Fail to run the container %v\n", err)
	//	}
	//
	//	logs, err := container.GetLogs()
	//	if err != nil {
	//		log.Fatalf("Fail to parse logs %v\n", err)
	//	}
	//
	//	log.Print(logs)
	//}

	//parameters := make(map[string]interface{})
	//datum := make(map[string]interface{})

	// Open file
	yfile, err := ioutil.ReadFile("examples/example2/workflow.yaml")

	if err != nil {

		log.Fatal(err)
	}

	arguments := make(map[string]string)

	arguments["git_repo"] = "https://github.com/Anthony-Jhoiro/sample_git.git"
	arguments["output"] = "test-wf-1"

	res, err := configParser.ParseWorkflowFile(yfile, arguments)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Res", res)

	res.Workflow.Run()

	//// Decode first part of the file
	//dec := yaml.NewDecoder(bytes.NewReader(yfile))
	//
	//// decode parameter template
	//err = dec.Decode(&parameters)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// TODO : fill parameters
	//arguments := make(map[string]string)
	//
	//arguments["repo_name"] = "https://github.com/Anthony-Jhoiro/sample_git.git"
	//arguments["output"] = "toto"
	//
	//// Decode main file
	//
	//params := struct {
	//	Parameters map[string]string
	//}{arguments}
	//
	//
	//tmpl, err := template.New("configParser").Parse( string(yfile))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//buf := &bytes.Buffer{}
	//
	//err = tmpl.Execute(buf, params)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// decode parser
	//
	//dec = yaml.NewDecoder(bytes.NewReader(buf.Bytes()))
	//
	//// decode parameter template
	//err = dec.Decode(&datum)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// decode parameter template
	//err = dec.Decode(&datum)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println(datum)
}
