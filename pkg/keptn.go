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
	OnTriggered(ce interface{}, context Context) (error, Context)

	OnFinished(context Context) interface{}

	GetTask() string

	GetData() interface{}
}

type KeptnOption func(*Keptn)

func WithHandler(handler TaskHandler) KeptnOption {
	return func(k *Keptn) {

		k.TaskRegistry.Add(keptnEventTypePrefix+handler.GetTask()+keptnTriggeredEventSuffix, TaskEntry{TaskHandler: handler})
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
		TaskRegistry:     NewTasksMap(),
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
	TaskRegistry     TaskRegistry
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
	if handler, ok := k.TaskRegistry.Contains(event.Type()); ok {
		data := handler.TaskHandler.GetData()
		if err := event.DataAs(&data); err != nil {
			k.handleErr(err)
		}
		if k.SendStartEvent {
			k.send(k.createStartedEventForTriggeredEvent(event))
		}
		err, newContext := handler.TaskHandler.OnTriggered(data, handler.Context)
		if err != nil {
			k.handleErr(err)
		}
		eventFinishedData := handler.TaskHandler.OnFinished(newContext)

		if k.SendFinishEvent {
			k.send(k.createFinishedEventForTriggeredEvent(event, eventFinishedData))
		}
	}
}

func (k Keptn) send(event cloudevents.Event) error {
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

func (k Keptn) createFinishedEventForTriggeredEvent(triggeredEvent cloudevents.Event, eventData interface{}) cloudevents.Event {
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
