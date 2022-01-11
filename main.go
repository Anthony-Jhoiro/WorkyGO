package main

import (
	"Workflow/configParser"
	"io/ioutil"
	"log"
)

func main() {
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
}
