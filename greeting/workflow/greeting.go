package main

import (
	"context"
	"encoding/json"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"learning.temporal/greeting"
	"log"
	"os"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	input := greeting.WorkflowInput{Name: os.Args[1], LanguageCode: os.Args[2]}

	options := client.StartWorkflowOptions{
		ID:                    "my-first-workflow-" + input.LanguageCode,
		TaskQueue:             greeting.TaskQueue,
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, greeting.GreetSomeone, input)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	var output greeting.WorkflowOutput
	err = we.Get(context.Background(), &output)
	if err != nil {
		log.Fatalln("Unable get workflow output", err)
	}
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalln("Unable to format result in JSON format", err)
	}
	log.Printf("Workflow result: %s\n", string(data))
}
