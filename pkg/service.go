package keptn

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Handler interface {
	Triggered(event cloudevents.Event)
	Type() string
}

type EventSender struct {
	//TODO: Add dependencies
}

func (es *EventSender) SendStartedEvent(event cloudevents.Event) error {
	fmt.Println("Sending .started event")
	//TODO: implement
	return nil
}

func (es *EventSender) SendFinishedEvent(event cloudevents.Event) error {
	fmt.Println("sending .finished event")
	//TODO: implement
	return nil
}

func NewKeptnService(ceClient cloudevents.Client, opts ...KeptnServiceOption) *KeptnService {
	k := &KeptnService{
		ceClient: ceClient,
		handlers: make(map[string]Handler),
	}
	for _, opt := range opts {
		opt(k)
	}
	return k
}

type KeptnService struct {
	ceClient cloudevents.Client
	handlers map[string]Handler
}

func (s KeptnService) Start() {
	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)
	s.ceClient.StartReceiver(ctx, s.gotEvent)
}

type KeptnServiceOption func(*KeptnService)

func WithHandler(h Handler) KeptnServiceOption {
	return func(ks *KeptnService) {
		ks.handlers[h.Type()] = h
	}
}

func (ks *KeptnService) gotEvent(event cloudevents.Event) {
	if h, ok := ks.handlers[event.Type()]; ok {
		h.Triggered(event)
	}
}
