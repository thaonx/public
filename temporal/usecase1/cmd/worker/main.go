package main

import (
	"log"
	"usecase1"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()


	for i := 0; i < 5; i++ {
		go func() {
			w := worker.New(c, "hello-world", worker.Options{
				// MaxConcurrentWorkflowTaskExecutionSize: ,
			})
			w.RegisterWorkflow(usecase1.WelcomeWorkflow)
			err = w.Run(worker.InterruptCh())
			if err != nil {
				log.Fatalln("Unable to start worker", err)
			}
		}()
	}

	done := make(chan bool, 1)
	<-done

	// numberOfWorker := 5
	// for i := 0; i < numberOfWorker; i++ {
	// 	go func() {

	// 	}()
	// }

}
