// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/warber/keptn-sdk-go/pkg/sdk"
	"sync"
)

// Ensure, that EventSenderMock does implement sdk.EventSender.
// If this is not the case, regenerate this file with moq.
var _ sdk.EventSender = &EventSenderMock{}

// EventSenderMock is a mock implementation of sdk.EventSender.
//
// 	func TestSomethingThatUsesEventSender(t *testing.T) {
//
// 		// make and configure a mocked sdk.EventSender
// 		mockedEventSender := &EventSenderMock{
// 			SendEventFunc: func(eventMoqParam event.Event) error {
// 				panic("mock out the SendEvent method")
// 			},
// 		}
//
// 		// use mockedEventSender in code that requires sdk.EventSender
// 		// and then make assertions.
//
// 	}
type EventSenderMock struct {
	// SendEventFunc mocks the SendEvent method.
	SendEventFunc func(eventMoqParam event.Event) error

	// calls tracks calls to the methods.
	calls struct {
		// SendEvent holds details about calls to the SendEvent method.
		SendEvent []struct {
			// EventMoqParam is the eventMoqParam argument value.
			EventMoqParam event.Event
		}
	}
	lockSendEvent sync.RWMutex
}

// SendEvent calls SendEventFunc.
func (mock *EventSenderMock) SendEvent(eventMoqParam event.Event) error {
	if mock.SendEventFunc == nil {
		panic("EventSenderMock.SendEventFunc: method is nil but EventSender.SendEvent was just called")
	}
	callInfo := struct {
		EventMoqParam event.Event
	}{
		EventMoqParam: eventMoqParam,
	}
	mock.lockSendEvent.Lock()
	mock.calls.SendEvent = append(mock.calls.SendEvent, callInfo)
	mock.lockSendEvent.Unlock()
	return mock.SendEventFunc(eventMoqParam)
}

// SendEventCalls gets all the calls that were made to SendEvent.
// Check the length with:
//     len(mockedEventSender.SendEventCalls())
func (mock *EventSenderMock) SendEventCalls() []struct {
	EventMoqParam event.Event
} {
	var calls []struct {
		EventMoqParam event.Event
	}
	mock.lockSendEvent.RLock()
	calls = mock.calls.SendEvent
	mock.lockSendEvent.RUnlock()
	return calls
}
