package run

import (
	"Workflow/configParser"
	"Workflow/logger"
	"Workflow/stepMapper"
	"Workflow/workflow"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"
)

func Run(filename string, arguments map[string]string, liveMode bool) {
	runNumber := strconv.FormatInt(time.Now().Unix(), 16)

	// Configure logger
	l := buildLogger(liveMode, runNumber)

	// Open logFile
	binaryTemplate, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	parsedWorkflow, err := configParser.ParseWorkflowFile(binaryTemplate, arguments)
	if err != nil {
		log.Fatal(err)
	}

	parsedWorkflow.SetLogger(l)
	parsedWorkflow.SetRunNumber(runNumber)

	wf, err := stepMapper.ParseWorkflowSteps(*parsedWorkflow)

	if err != nil {
		log.Fatal(err)
	}

	var wfError error

	if liveMode {
		signals := make(chan bool)

		go liveLogWorkflow(signals, wf)
		wfError = wf.Run(parsedWorkflow)

		signals <- true
	} else {
		wfError = wf.Run(parsedWorkflow)
	}

	cleanWorkflowSteps(wf)

	if wfError != nil {
		log.Fatal(wfError)
	}

}

func buildLogger(liveMode bool, runNumber string) logger.Logger {
	if !liveMode {
		return logger.New("", os.Stdout)
	} else {

		historyPath := path.Join("./history", fmt.Sprintf("run-%s", runNumber))

		err := os.MkdirAll(historyPath, os.ModePerm)
		if err != nil {
			log.Fatalf("can not create run history directory : %v", err)
		}

		file, err := os.Create(path.Join(historyPath, "run.log"))
		return logger.New("", file)
	}
}

func cleanWorkflowSteps(wf *workflow.Workflow) {
	for _, step := range wf.Steps {
		step.Clean()
	}
}

func liveLogWorkflow(ch chan bool, workflow *workflow.Workflow) {
	loop := true
	for loop {
		select {
		case <-ch:
			loop = false
		default:

		}
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
		workflow.Print()
		time.Sleep(1 * time.Second)
	}
}
