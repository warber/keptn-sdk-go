# keptn-sdk-go (experimental)

## Example

```go
package main

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/warber/keptn-sdk-go/pkg/sdk"
	"log"
)

func main() {
	// 1. create a HTTP client
	httpClient := sdk.GetHTTPClient(cloudevents.WithPath("/"), cloudevents.WithPort(8080))

	// 2. create keptn and register our handlers
	myKeptn := sdk.NewKeptn(httpClient, "my-service", sdk.WithHandler(GreetingsHandler{}, "sh.keptn.event.greeting.triggered"))

	// 3. start
	myKeptn.Start()
}

// GreetingsHandler is our custom handler to handle greetings.triggered events
type GreetingsHandler struct {
}

// Execute is called whenever a "greetings.triggered" event was received
func (m GreetingsHandler) Execute(ce interface{}, context sdk.Context) (error, sdk.Context) {
	greetingTriggeredData := ce.(*GreetingTriggeredData)

	// [...]
	// do stuff with your data
	// [...]

	// when we are done we set the finished data
	context.SetFinishedData(GreetingFinishedData{GreetMessage: greetingTriggeredData.GreetMessage})
	return nil, context
}

// GetData returns the type of data we want to process
func (m GreetingsHandler) GetData() interface{} {
	return &GreetingTriggeredData{}
}

// GreetingTriggeredData is our custom event data we want to process
type GreetingTriggeredData struct {
	keptnv2.EventData
	GreetMessage string `json:"greetMessage"`
}

// GreetingFInishedData is our custom event data we want to send out when our task processing is finished
type GreetingFinishedData struct {
	keptnv2.EventData
	GreetMessage string `json:"greetMessage"`
}

```