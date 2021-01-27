package example

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptn "github.com/warber/keptn-sdk/pkg"
)

type MyHandler struct {
	keptn.EventSender
}

func (h *MyHandler) Triggered(event cloudevents.Event) {
	fmt.Println("Handle Event")
	h.SendStartedEvent(event)

	//LOGIC

	h.SendFinishedEvent(event)
}

func (h *MyHandler) Type() string {
	return "keptn.sh.event.task.triggered"
}
