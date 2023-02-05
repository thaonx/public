package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"usecase1"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	mux := http.NewServeMux()
	count := 1

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
				// find out exactly what the error was and set err
				var err error
				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("Unknown panic")
				}
				if err != nil {
					// sendMeMail(err)
					fmt.Println("sendMeMail")
				}
			}
		}()
		payload := struct {
			Name string `json:"name"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			panic(err)
		}

		count++
		workflowOptions := client.StartWorkflowOptions{
			ID:        "helloWorkflowID-" + strconv.Itoa(count),
			TaskQueue: "hello-world",
		}

		welcomeInput := usecase1.WelcomeInput{
			Name: payload.Name + " " + strconv.Itoa(count),
		}
		we, err := c.ExecuteWorkflow(r.Context(), workflowOptions, usecase1.WelcomeWorkflow, welcomeInput)
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

		// Synchronously wait for the workflow completion.
		var welcomeOutput usecase1.WelcomeOutput
		err = we.Get(r.Context(), &welcomeOutput)
		if err != nil {
			log.Fatalln("Unable get workflow result", err)
		}
		if err := json.NewEncoder(w).Encode(welcomeOutput); err != nil {
			panic(err)
		}
	})
	if err := http.ListenAndServe(":9000", mux); err != nil {
		panic(err)
	}
}
