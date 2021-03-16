package main

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/warber/keptn-sdk-go/pkg/sdk"
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
	httpClient := sdk.GetHTTPClient(cloudevents.WithPath("/"), cloudevents.WithPort(8080))

	// 2. create keptn
	myKeptn := sdk.NewKeptn(httpClient, "my-service", sdk.WithHandler(DeploymentHandler{}))

	// 3. start
	myKeptn.Start()

}

// DeploymentHandler handles "sh.keptn.event.deployment.triggered" tasks
type DeploymentHandler struct {
}

func (m DeploymentHandler) OnTriggered(ce interface{}, context sdk.Context) (error, sdk.Context) {
	greetTriggeredData := ce.(*GreetTriggeredData)
	log.Println("Got GreetTriggered Event: " + greetTriggeredData.GreetMessage)
	<-time.After(5 * time.Second)
	context.SetFinishedData(GreetFinishedData{GreetMessage: greetTriggeredData.GreetMessage})
	return nil, context
}

func (m DeploymentHandler) GetTask() string {
	return "greet"
}

func (m DeploymentHandler) GetData() interface{} {
	return &GreetTriggeredData{}
}

type GreetTriggeredData struct {
	keptnv2.EventData
	GreetMessage string `json:"greetMessage"`
}

type GreetFinishedData struct {
	keptnv2.EventData
	GreetMessage string `json:"greetMessage"`
}
