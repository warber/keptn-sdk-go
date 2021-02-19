package keptn

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"log"
)

type OnTriggeredFunc func(ce cloudevents.Event) error

type KeptnOption func(*Keptn)

func OnTask(taskName string) KeptnOption {
	return func(k *Keptn) {
		k.OnTask = taskName
	}
}

func OnTriggered(onTriggeredFunc OnTriggeredFunc) KeptnOption {
	return func(k *Keptn) {
		k.OnTriggeredFunc = onTriggeredFunc
	}
}

func SendStartEvent(sendStartEvent bool) KeptnOption {
	return func(k *Keptn) {
		k.SendStartEvent = sendStartEvent
	}
}

func SendFinishEvent(sendFinishEvent bool) KeptnOption {
	return func(k *Keptn) {
		k.SendFInishEvent = sendFinishEvent
	}
}

func NewKeptn(ceClient cloudevents.Client, opts ...KeptnOption) *Keptn {
	keptn := &Keptn{
		CloudEventClient: ceClient,
	}
	for _, opt := range opts {
		opt(keptn)
	}
	return keptn
}

type Keptn struct {
	CloudEventClient cloudevents.Client
	OnTask           string
	OnTriggeredFunc  OnTriggeredFunc
	SendStartEvent   bool
	SendFInishEvent  bool
}

func (k Keptn) Start() {
	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)
	k.CloudEventClient.StartReceiver(ctx, k.gotEvent)
}

func (k Keptn) gotEvent(event cloudevents.Event) {
	k.sendTaskFinishedEvent(event)

	err := k.OnTriggeredFunc(event)
	if err != nil {
		k.handleErr(err)
	}

	k.sendTaskFinishedEvent(event)

}

func (k Keptn) sendTaskStartedEvent(ce cloudevents.Event) error {
	log.Println("Sending Task Started Event")
	return nil
}

func (k Keptn) sendTaskFinishedEvent(ce cloudevents.Event) error {
	log.Println("Sending Task Finished Event")
	return nil
}

func (k Keptn) handleErr(err error) {
	log.Println("Handling Error")
	//TODO implement
}
