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
	"strings"
	"time"
)

func getOutputParser(output map[string]map[string]string) func(string, string) (string, error) {
	return func(stepName, varName string) (string, error) {
		step, ok := output[stepName]
		if !ok {
			return "", fmt.Errorf("step %s not found", stepName)
		}

		value, ok := step[varName]
		if !ok {
			return "", fmt.Errorf("output value %s not found", varName)
		}

		return value, nil
	}
}

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

	var output map[string]map[string]string

	if liveMode {
		signals := make(chan bool)

		go liveLogWorkflow(signals, wf)
		output, wfError = wf.Run(parsedWorkflow)

		signals <- true
	} else {
		output, wfError = wf.Run(parsedWorkflow)
	}

	cleanWorkflowSteps(wf)

	parser := getOutputParser(output)

	ov := make(map[string]string)

	for outputName, key := range parsedWorkflow.Output {
		slicedKey := strings.Split(key, ".")

		if len(slicedKey) != 2 {
			log.Fatalf("output key %s is malformed", key)
		}
		value, err := parser(slicedKey[0], slicedKey[1])
		if err != nil {
			log.Fatalf("no value founded for key %s : %v", key, err)

		}
		ov[outputName] = value
	}

	for key, value := range ov {
		_ = l.Print(fmt.Sprintf("%s => %s", key, value))
	}

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
