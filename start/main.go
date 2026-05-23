package main

import (
	"context"
	"fmt"
	"log"

	"go.temporal.io/sdk/client"

	"temporal-hand-on/101/app"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal("Unable to create Temporal client", err)
	}
	defer c.Close()

	input := app.ClaimRequest{
		ClaimID:  "CLM-001",
		PolicyID: "POL-123",
		Amount:   15000,
	}

	we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:        "claim-001",
		TaskQueue: "car-insurance-claims",
	}, app.ProcessClaim, input)
	if err != nil {
		log.Fatal("Unable to start workflow", err)
	}

	fmt.Printf("Started workflow: WorkflowID=%s RunID=%s\n", we.GetID(), we.GetRunID())

	var result app.ClaimResult
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatal("Workflow failed", err)
	}

	fmt.Printf("Result: ClaimID=%s Status=%s Reason=%s\n", result.ClaimID, result.Status, result.Reason)
}
