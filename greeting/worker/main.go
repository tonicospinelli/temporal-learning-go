package main

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"learning.temporal/greeting"
	"log"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "greeting-tasks", worker.Options{})

	w.RegisterWorkflow(greeting.GreetSomeone)
	w.RegisterActivity(greeting.GreetInSpanish)
	w.RegisterActivity(greeting.FarewellInSpanish)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
