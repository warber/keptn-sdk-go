package keptn

import (
	"context"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	httpprotocol "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"time"
)

type Keptn struct {
	ceClient       cloudevents.Client
	eventBrokerURL string
}

func (*Keptn) SendStartedEvent(event cloudevents.Event) {
	fmt.Println("SEND STARTED EVENT")
}

func (*Keptn) SendFinishedEvent(event cloudevents.Event) {
	fmt.Println("SEND FINISHED EVENT")
}

func (*Keptn) SendEvent(i interface{}) {
	fmt.Println("SEND EVENT")
}

func (k *Keptn) sendCloudEvent(event cloudevents.Event) error {
	ctx := cloudevents.ContextWithTarget(context.Background(), k.eventBrokerURL)
	ctx = cloudevents.WithEncodingStructured(ctx)
	var result protocol.Result
	for i := 0; i <= 3; i++ {
		result = k.ceClient.Send(ctx, event)
		httpResult, ok := result.(*httpprotocol.Result)
		if ok {
			if httpResult.StatusCode >= 200 && httpResult.StatusCode < 300 {
				return nil
			} else {
				<-time.After(keptn.GetExpBackoffTime(i + 1))
			}
		} else if cloudevents.IsUndelivered(result) {
			<-time.After(keptn.GetExpBackoffTime(i + 1))
		} else {
			return nil
		}
	}
	return errors.New("Failed to send cloudevent: " + result.Error())
}

type KeptnApiContext struct {
	keptn   *Keptn
	context context.Context
}

func (c *KeptnApiContext) GetKeptn() *Keptn {
	return c.keptn
}

type EventHandler func(ctx KeptnApiContext, event cloudevents.Event) error

type EventRecorder struct {
	handlerMap map[string]EventHandler
}

func nilHandler(ctx KeptnApiContext, event cloudevents.Event) error {
	fmt.Printf("Ignoring event %s\n", event.Type())
	return nil
}

func (er *EventRecorder) Add(eventType string, handler EventHandler) {
	er.handlerMap[eventType] = handler
}

func (er *EventRecorder) Get(eventType string) EventHandler {
	if handler, ok := er.handlerMap[eventType]; ok {
		return handler
	}
	return nilHandler
}

func NewKeptnController(ceClient cloudevents.Client) *KeptnController {
	return &KeptnController{
		eventRecorder: EventRecorder{
			handlerMap: make(map[string]EventHandler),
		},
		ceClient: ceClient,
	}
}

type KeptnController struct {
	eventRecorder EventRecorder
	ceClient      cloudevents.Client
}

func (k *KeptnController) Register(eventType string, handler EventHandler) error {
	k.eventRecorder.Add(eventType, handler)
	return nil
}

func (k *KeptnController) Start() {
	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)
	k.ceClient.StartReceiver(ctx, k.gotEvent)
}

func (k *KeptnController) gotEvent(ctx context.Context, event cloudevents.Event) {
	go k.eventRecorder.Get(event.Type())(KeptnApiContext{
		context: ctx,
		keptn:   &Keptn{ceClient: k.ceClient},
	}, event)
}
