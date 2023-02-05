package main

import (
	"context"
	"log"
	"usecase1"

	"go.temporal.io/sdk/client"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		// ID:        "hello_world_workflowID",
		TaskQueue: "hello-world",
	}

	welcomeIn := usecase1.WelcomeInput{
		Name: "Thao Nguyen",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, usecase1.WelcomeWorkflow, welcomeIn)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var welcomeOutput usecase1.WelcomeOutput
	err = we.Get(context.Background(), &welcomeOutput)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Printf("Workflow result: %#v \n", welcomeOutput)
}
