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

func (m DeploymentHandler) OnTriggered(ce interface{}, context keptn.Context) (error, keptn.Context) {
	greetTriggeredData := ce.(*GreetTriggeredData)
	log.Println("Got GreetTriggered Event: " + greetTriggeredData.GreetMessage)
	<-time.After(5 * time.Second)
	context.SetFinishedData(GreetFinishedData{GreetMessage: greetTriggeredData.GreetMessage})
	return nil, context
}

func (m DeploymentHandler) OnFinished(context keptn.Context) interface{} {
	log.Println("Executing OnFinish Logic")
	log.Println(context.FinishedData)
	return context.FinishedData
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
