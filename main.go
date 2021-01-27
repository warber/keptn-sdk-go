package main

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptn "github.com/warber/keptn-sdk/pkg"
	"log"
)

func main() {

	keptn := keptn.NewKeptnController(createClient())
	keptn.Register("keptn.sh.event.task.triggered", handleTaskTriggered)
	keptn.Start()
}

func handleTaskTriggered(ctx keptn.KeptnApiContext, ce cloudevents.Event) error {
	fmt.Println("Handling keptn.sh.event.task.triggered event")
	ctx.GetKeptn().SendStartedEvent(ce)

	// EXECUTE YOUR LOGIC

	ctx.GetKeptn().SendFinishedEvent(ce)

	return nil
}

func createClient() cloudevents.Client {
	p, err := cloudevents.NewHTTP(cloudevents.WithPath("/"), cloudevents.WithPort(8080))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, _ := cloudevents.NewClient(p)
	return c
}
