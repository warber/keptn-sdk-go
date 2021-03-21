package sdk_test

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/warber/keptn-sdk-go/pkg/sdk"
	"github.com/warber/keptn-sdk-go/pkg/sdk/fake"
	"testing"
)

func Test_Keptn(t *testing.T) {

	taskHandler := FakeTaskHandler{}
	taskContext := sdk.Context{}

	taskEntry := sdk.TaskEntry{
		TaskHandler: taskHandler,
		Context:     taskContext,
	}

	taskEntries := map[string]sdk.TaskEntry{"sh.keptn.event.faketask.triggered": taskEntry}

	eventReceiver := &fake.TestReceiver{}
	eventSender := &fake.EventSenderMock{}
	taskRegistry := sdk.TaskRegistry{
		Entries: taskEntries,
	}

	keptn := sdk.Keptn{
		EventSender:     eventSender,
		EventReceiver:   eventReceiver,
		TaskRegistry:    taskRegistry,
		SendStartEvent:  false,
		SendFinishEvent: false,
	}

	keptn.Start()
	eventReceiver.NewEvent(newTestTaskTriggeredEvent())
	//TODO: assert eventSender
}

func newTestTaskTriggeredEvent() cloudevents.Event {
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType("sh.keptn.event.faketask.triggered")
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(sdk.KeptnContextCEExtension, "keptncontext")
	c.SetExtension(sdk.TriggeredIDCEExtension, "ID")
	c.SetSource("unittest")
	c.SetData(cloudevents.ApplicationJSON, FakeTaskData{})
	return c
}

type FakeTaskData struct {
}

type FakeTaskHandler struct {
}

func (f FakeTaskHandler) Execute(ce interface{}, context sdk.Context) (error, sdk.Context) {
	fmt.Println("FakeTaskHandler::Execute() called")
	return nil, context
}

func (f FakeTaskHandler) GetData() interface{} {
	return FakeTaskData{}
}
