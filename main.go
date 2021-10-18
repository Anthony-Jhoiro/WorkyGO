package main

import (
	"Workflow/workflow/steps"
	"log"
)

func main() {
	stepA := steps.CreateStep("A")
	stepB := steps.CreateStep("B")
	stepC := steps.CreateStep("C")
	stepD := steps.CreateStep("D")

	stepB.AddRequirement(stepA)
	stepC.AddRequirement(stepA)
	stepD.AddRequirement(stepB)
	stepD.AddRequirement(stepC)

	_, err := stepA.Execute(nil)

	if err != nil {
		log.Println("Something failed")
	}

}
