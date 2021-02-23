package main

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	keptn "github.com/warber/keptn-sdk/pkg"
	"log"
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

	myKeptn := keptn.NewKeptn(createClient(), "my-service", keptn.WithHandler(MyHandler{}))
	myKeptn.Start()

}

type MyHandler struct {
}

func (m MyHandler) OnTriggered(ce cloudevents.Event) error {
	log.Println("Executing Business Logic")
	return nil
}

func (m MyHandler) OnFinished() keptnv2.EventData {
	log.Println("Executing OnFinish Logic")
	return keptnv2.EventData{}
}

func (m MyHandler) GetTask() string {
	return "deployment"
}
