// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	"github.com/warber/keptn-sdk-go/pkg/sdk"
	"sync"
)

// Ensure, that TaskHandlerMock does implement sdk.TaskHandler.
// If this is not the case, regenerate this file with moq.
var _ sdk.TaskHandler = &TaskHandlerMock{}

// TaskHandlerMock is a mock implementation of sdk.TaskHandler.
//
// 	func TestSomethingThatUsesTaskHandler(t *testing.T) {
//
// 		// make and configure a mocked sdk.TaskHandler
// 		mockedTaskHandler := &TaskHandlerMock{
// 			ExecuteFunc: func(ce interface{}, context sdk.Context) (error, sdk.Context) {
// 				panic("mock out the Execute method")
// 			},
// 			GetDataFunc: func() interface{} {
// 				panic("mock out the GetData method")
// 			},
// 		}
//
// 		// use mockedTaskHandler in code that requires sdk.TaskHandler
// 		// and then make assertions.
//
// 	}
type TaskHandlerMock struct {
	// ExecuteFunc mocks the Execute method.
	ExecuteFunc func(ce interface{}, context sdk.Context) (error, sdk.Context)

	// GetDataFunc mocks the GetData method.
	GetDataFunc func() interface{}

	// calls tracks calls to the methods.
	calls struct {
		// Execute holds details about calls to the Execute method.
		Execute []struct {
			// Ce is the ce argument value.
			Ce interface{}
			// Context is the context argument value.
			Context sdk.Context
		}
		// GetData holds details about calls to the GetData method.
		GetData []struct {
		}
	}
	lockExecute sync.RWMutex
	lockGetData sync.RWMutex
}

// Execute calls ExecuteFunc.
func (mock *TaskHandlerMock) Execute(ce interface{}, context sdk.Context) (error, sdk.Context) {
	if mock.ExecuteFunc == nil {
		panic("TaskHandlerMock.ExecuteFunc: method is nil but TaskHandler.Execute was just called")
	}
	callInfo := struct {
		Ce      interface{}
		Context sdk.Context
	}{
		Ce:      ce,
		Context: context,
	}
	mock.lockExecute.Lock()
	mock.calls.Execute = append(mock.calls.Execute, callInfo)
	mock.lockExecute.Unlock()
	return mock.ExecuteFunc(ce, context)
}

// ExecuteCalls gets all the calls that were made to Execute.
// Check the length with:
//     len(mockedTaskHandler.ExecuteCalls())
func (mock *TaskHandlerMock) ExecuteCalls() []struct {
	Ce      interface{}
	Context sdk.Context
} {
	var calls []struct {
		Ce      interface{}
		Context sdk.Context
	}
	mock.lockExecute.RLock()
	calls = mock.calls.Execute
	mock.lockExecute.RUnlock()
	return calls
}

// GetData calls GetDataFunc.
func (mock *TaskHandlerMock) GetData() interface{} {
	if mock.GetDataFunc == nil {
		panic("TaskHandlerMock.GetDataFunc: method is nil but TaskHandler.GetData was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetData.Lock()
	mock.calls.GetData = append(mock.calls.GetData, callInfo)
	mock.lockGetData.Unlock()
	return mock.GetDataFunc()
}

// GetDataCalls gets all the calls that were made to GetData.
// Check the length with:
//     len(mockedTaskHandler.GetDataCalls())
func (mock *TaskHandlerMock) GetDataCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetData.RLock()
	calls = mock.calls.GetData
	mock.lockGetData.RUnlock()
	return calls
}
