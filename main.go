package main

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	keptn "github.com/warber/keptn-sdk/pkg"
	"log"
	"time"
)

func createClient() cloudevents.Client {
	p, err := cloudevents.NewHTTP(cloudevents.WithPath("/"), cloudevents.WithPort(8080))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, _ := cloudevents.NewClient(p)
	return c
}

func main() {
	// 1. create http client
	httpClient := keptn.GetHTTPClient(cloudevents.WithPath("/"), cloudevents.WithPort(8080))

	// 2. create keptn
	myKeptn := keptn.NewKeptn(httpClient, "my-service", keptn.WithHandler(DeploymentHandler{}))

	// 3. start
	myKeptn.Start()

}

// DeploymentHandler handles "sh.keptn.event.deployment.triggered" tasks
type DeploymentHandler struct {
}

func (m DeploymentHandler) OnTriggered(ce cloudevents.Event) error {
	log.Println("Executing Business Logic")
	<-time.After(3 * time.Second)
	return nil
}

func (m DeploymentHandler) OnFinished() keptnv2.EventData {
	log.Println("Executing OnFinish Logic")
	return keptnv2.EventData{}
}

func (m DeploymentHandler) GetTask() string {
	return "greet"
}
