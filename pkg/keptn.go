package keptn

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"log"
	"strings"
)

type Handler interface {
	OnTriggered(ce cloudevents.Event) error
	OnFinished() keptnv2.EventData
	GetTask() string
}

type OnTriggeredFunc func(ce cloudevents.Event) error

type KeptnOption func(*Keptn)

func WithHandler(handler Handler) KeptnOption {
	return func(k *Keptn) {

		k.Handlers[keptnEventTypePrefix+handler.GetTask()+keptnTriggeredEventSuffix] = handler
	}
}

func SendStartEvent(sendStartEvent bool) KeptnOption {
	return func(k *Keptn) {
		k.SendStartEvent = sendStartEvent
	}
}

func SendFinishEvent(sendFinishEvent bool) KeptnOption {
	return func(k *Keptn) {
		k.SendFinishEvent = sendFinishEvent
	}
}

func NewKeptn(ceClient cloudevents.Client, source string, opts ...KeptnOption) *Keptn {
	keptn := &Keptn{
		CloudEventClient: ceClient,
		Source:           source,
		Handlers:         make(map[string]Handler),
		SendStartEvent:   true,
		SendFinishEvent:  true,
	}
	for _, opt := range opts {
		opt(keptn)
	}
	return keptn
}

type Keptn struct {
	CloudEventClient cloudevents.Client
	Source           string
	Handlers         map[string]Handler
	SendStartEvent   bool
	SendFinishEvent  bool
}

func (k Keptn) Start() {
	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)
	k.CloudEventClient.StartReceiver(ctx, k.gotEvent)
}

func (k Keptn) gotEvent(event cloudevents.Event) {
	if handler, ok := k.Handlers[event.Type()]; ok {
		k.sendTaskStartedEvent(event)
		if err := handler.OnTriggered(event); err != nil {
			k.handleErr(err)
		}
		eventFinishedData := handler.OnFinished()
		finishedEvent := k.createFinishedEventForTriggeredEvent(event, eventFinishedData)
		k.sendTaskFinishedEvent(finishedEvent)
	}

}

func (k Keptn) sendTaskStartedEvent(triggeredEvent cloudevents.Event) error {
	log.Println("Sending Task Started Event")
	//TODO implement
	return nil
}

func (k Keptn) sendTaskFinishedEvent(ce cloudevents.Event) error {
	log.Println("Sending Task Finished Event")
	//TODO implement
	return nil
}

func (k Keptn) handleErr(err error) {
	log.Println("Handling Error")
	//TODO implement
}

func (k Keptn) createFinishedEventForTriggeredEvent(triggeredEvent cloudevents.Event, eventData keptnv2.EventData) cloudevents.Event {
	finishedEventType := strings.Trim(triggeredEvent.Type(), ".triggered") + ".finished"
	keptnContext, _ := triggeredEvent.Context.GetExtension(keptnContextCEExtension)
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(finishedEventType)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(keptnContextCEExtension, keptnContext)
	c.SetExtension(triggeredIDCEExtension, triggeredEvent.ID())
	c.SetSource(k.Source)
	c.SetData(cloudevents.ApplicationJSON, keptnv2.EventData{})
	return c
}
