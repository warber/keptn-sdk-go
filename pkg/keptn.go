package keptn

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"log"
	"strings"
)

type KeptnEventData struct {
	Project string
	Stage   string
	Service string
}

type TaskHandler interface {

	// OnTriggered is called when a event was received that shall triggeres the TaskHandler
	// Implement the business logic to process the task here
	OnTriggered(ce interface{}) error

	// OnFinished is called after the logic in OnTriggered was executed, i.e., after the task was processed.
	// It shall return a keptnv2.EventData object that holds the data to be contained in the
	// .finished event eventually sent by Keptn
	OnFinished() keptnv2.EventData

	// GetTask returns the task name for which the TaskHandler is responsible
	GetTask() string

	GetData() interface{}
}

type KeptnOption func(*Keptn)

func WithHandler(handler TaskHandler) KeptnOption {
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
		EventSender:      NewHTTPEventSender(ceClient),
		CloudEventClient: ceClient,
		Source:           source,
		Handlers:         make(map[string]TaskHandler),
		SendStartEvent:   true,
		SendFinishEvent:  true,
	}
	for _, opt := range opts {
		opt(keptn)
	}
	return keptn
}

type Keptn struct {
	EventSender      EventSender
	CloudEventClient cloudevents.Client
	Source           string
	Handlers         map[string]TaskHandler
	SendStartEvent   bool
	SendFinishEvent  bool
}

func (k Keptn) Start() {
	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)
	err := k.CloudEventClient.StartReceiver(ctx, k.gotEvent)
	_ = err
}

func (k Keptn) gotEvent(event cloudevents.Event) {

	if handler, ok := k.Handlers[event.Type()]; ok {

		data := handler.GetData()
		if err := event.DataAs(&data); err != nil {
			k.handleErr(err)
		}
		if k.SendStartEvent {
			k.send(k.createStartedEventForTriggeredEvent(event))
		}
		if err := handler.OnTriggered(data); err != nil {
			k.handleErr(err)
		}
		eventFinishedData := handler.OnFinished()

		if k.SendFinishEvent {
			k.send(k.createFinishedEventForTriggeredEvent(event, eventFinishedData))
		}
	}

}

func (k Keptn) send(event cloudevents.Event) error {
	log.Println("Sending Task Started Event")
	if err := k.EventSender.SendEvent(event); err != nil {
		log.Println("Error sending .started event")
	}
	return nil
}

func (k Keptn) handleErr(err error) {
	log.Println("Handling Error")
	//TODO SEND EVENT
}

func (k Keptn) createStartedEventForTriggeredEvent(triggeredEvent cloudevents.Event) cloudevents.Event {
	fmt.Println(triggeredEvent.Type())

	startedEventType := strings.TrimSuffix(triggeredEvent.Type(), ".triggered") + ".started"
	keptnContext, _ := triggeredEvent.Context.GetExtension(keptnContextCEExtension)
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(startedEventType)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(keptnContextCEExtension, keptnContext)
	c.SetExtension(triggeredIDCEExtension, triggeredEvent.ID())
	c.SetSource(k.Source)
	c.SetData(cloudevents.ApplicationJSON, keptnv2.EventData{})
	return c
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
	c.SetData(cloudevents.ApplicationJSON, eventData)
	return c
}
