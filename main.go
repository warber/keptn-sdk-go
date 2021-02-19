package main

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
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

	myKeptn := keptn.NewKeptn(createClient(), keptn.OnTriggered(handleEvent))
	myKeptn.Start()

}

func handleEvent(ce cloudevents.Event) error {
	log.Println("Handling event")
	return nil
}
