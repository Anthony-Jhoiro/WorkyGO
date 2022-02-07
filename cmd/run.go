/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"Workflow/run"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
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
		run.Run(args[0], arguments, LiveMode)
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
