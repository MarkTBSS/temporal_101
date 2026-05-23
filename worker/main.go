package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"temporal-hand-on/101/app"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal("Unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, "car-insurance-claims", worker.Options{})

	w.RegisterWorkflow(app.ProcessClaim)
	w.RegisterActivity(app.VerifyPolicy)
	w.RegisterActivity(app.ApproveClaim)
	w.RegisterActivity(app.NotifyCustomer)

	log.Println("Worker started, listening on task queue: car-insurance-claims")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatal("Unable to start worker", err)
	}
}
