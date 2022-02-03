package main

import (
	"Workflow/configParser"
	"Workflow/logger"
	"Workflow/stepMapper"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

func main() {

	runNumber := strconv.FormatInt(time.Now().Unix(), 16)

	historyPath := path.Join("./history", fmt.Sprintf("run-%s", runNumber))

	err := os.MkdirAll(historyPath, os.ModePerm)
	if err != nil {
		log.Fatalf("can not create run history directory : %v", err)
	}

	file, err := os.Create(path.Join(historyPath, "run.log"))
	defer file.Close()

	l := logger.New("", file)

	// Open logFile
	yfile, err := ioutil.ReadFile("examples/example2/workflow.yaml")

	if err != nil {

		log.Fatal(err)
	}

	arguments := make(map[string]string)

	arguments["git_repo"] = "https://github.com/Anthony-Jhoiro/sample_git.git"
	arguments["output"] = "test-wf-1"

	parsedWorkflow, err := configParser.ParseWorkflowFile(yfile, arguments)
	if err != nil {
		log.Fatal(err)
	}

	parsedWorkflow.SetLogger(l)
	parsedWorkflow.SetRunNumber(runNumber)

	workflow, err := stepMapper.ParseWorkflowSteps(*parsedWorkflow)

	if err != nil {
		log.Fatal(err)
	}

	workflow.Run(parsedWorkflow)

	for _, step := range workflow.Steps {
		step.Clean()
	}
}
