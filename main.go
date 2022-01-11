package main

import (
	"Workflow/configParser"
	"Workflow/logger"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

func main() {

	err := logger.LOG.Init(logger.Context{RunName: strconv.FormatInt(time.Now().Unix(), 10)})
	if err != nil {
		log.Fatalf("Fail to initialize logger : %v", err)
	}

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

	res.Workflow.Run()
}
