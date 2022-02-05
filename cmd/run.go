/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"Workflow/configParser"
	"Workflow/logger"
	"Workflow/stepMapper"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var LiveMode bool
var CLIArguments []string
var JsonArgumentsFilename string
var arguments map[string]string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start a workflow",
	Long:  `Start a workflow with the given parameters.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Require a positional argument for the filename
		if len(args) != 1 {
			return errors.New("at least one argument is required for the workflow file")
		}

		arguments = make(map[string]string)

		if len(JsonArgumentsFilename) != 0 {
			fileContent, err := ioutil.ReadFile(JsonArgumentsFilename)
			if err != nil {
				return fmt.Errorf("fail to parse arguments file : %v", err)
			}
			err = json.Unmarshal(fileContent, &arguments)
			if err != nil {
				return fmt.Errorf("fail to parse arguments file : %v", err)
			}
		}

		// Check arguments format
		for index, argument := range CLIArguments {
			if !strings.Contains(argument, "=") {
				return fmt.Errorf("error in argument %d (%s) : invalid format ", index, argument)
			}
			slicedString := strings.SplitN(argument, "=", 2)
			arguments[slicedString[0]] = slicedString[1]
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		runNumber := strconv.FormatInt(time.Now().Unix(), 16)

		// Configure logger
		var l logger.Logger

		if !LiveMode {
			l = logger.New("", os.Stdout)
		} else {

			historyPath := path.Join("./history", fmt.Sprintf("run-%s", runNumber))

			err := os.MkdirAll(historyPath, os.ModePerm)
			if err != nil {
				log.Fatalf("can not create run history directory : %v", err)
			}

			file, err := os.Create(path.Join(historyPath, "run.log"))
			defer file.Close()
			l = logger.New("", file)
		}

		// Open logFile
		yfile, err := ioutil.ReadFile(args[0])

		if err != nil {

			log.Fatal(err)
		}

		//arguments := make(map[string]string)

		for _, argument := range CLIArguments {
			slicedString := strings.SplitN(argument, "=", 2)
			arguments[slicedString[0]] = slicedString[1]
		}

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

		if LiveMode {
			signals := make(chan bool)

			go func() {
				loop := true
				for loop {
					select {
					case <-signals:
						loop = false
					default:

					}
					cmd := exec.Command("clear") //Linux example, its tested
					cmd.Stdout = os.Stdout
					_ = cmd.Run()
					workflow.Print()
					time.Sleep(1 * time.Second)
				}

			}()
			workflow.Run(parsedWorkflow)

			signals <- true
		} else {
			workflow.Run(parsedWorkflow)
		}

		for _, step := range workflow.Steps {
			step.Clean()
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().BoolVarP(&LiveMode, "live", "l", false, "Use live mode")
	runCmd.PersistentFlags().StringArrayVarP(&CLIArguments, "arg", "a", []string{}, "Your arguments formatted as key=value (you can have multiple arguments by using multiple time this argument).")
	runCmd.PersistentFlags().StringVarP(&JsonArgumentsFilename, "from-json", "f", "", "A Json file that contains your arguments. The file must be formatted as {\"key\": \"value\"}. Only strings are supported as values")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
