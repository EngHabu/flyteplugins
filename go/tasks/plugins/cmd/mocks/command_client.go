// Code generated by mockery v1.0.1. DO NOT EDIT.

package mocks

import (
	context "context"

	command "github.com/lyft/flyteplugins/go/tasks/plugins/cmd"

	mock "github.com/stretchr/testify/mock"
)

// CommandClient is an autogenerated mock type for the CommandClient type
type CommandClient struct {
	mock.Mock
}

type CommandClient_ExecuteCommand struct {
	*mock.Call
}

func (_m CommandClient_ExecuteCommand) Return(_a0 interface{}, _a1 error) *CommandClient_ExecuteCommand {
	return &CommandClient_ExecuteCommand{Call: _m.Call.Return(_a0, _a1)}
}

func (_m *CommandClient) OnExecuteCommand(ctx context.Context, commandStr string, extraArgs interface{}) *CommandClient_ExecuteCommand {
	c := _m.On("ExecuteCommand", ctx, commandStr, extraArgs)
	return &CommandClient_ExecuteCommand{Call: c}
}

func (_m *CommandClient) OnExecuteCommandMatch(matchers ...interface{}) *CommandClient_ExecuteCommand {
	c := _m.On("ExecuteCommand", matchers...)
	return &CommandClient_ExecuteCommand{Call: c}
}

// ExecuteCommand provides a mock function with given fields: ctx, commandStr, extraArgs
func (_m *CommandClient) ExecuteCommand(ctx context.Context, commandStr string, extraArgs interface{}) (interface{}, error) {
	ret := _m.Called(ctx, commandStr, extraArgs)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) interface{}); ok {
		r0 = rf(ctx, commandStr, extraArgs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, interface{}) error); ok {
		r1 = rf(ctx, commandStr, extraArgs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type CommandClient_GetCommandStatus struct {
	*mock.Call
}

func (_m CommandClient_GetCommandStatus) Return(_a0 command.CommandStatus, _a1 error) *CommandClient_GetCommandStatus {
	return &CommandClient_GetCommandStatus{Call: _m.Call.Return(_a0, _a1)}
}

func (_m *CommandClient) OnGetCommandStatus(ctx context.Context, commandID string) *CommandClient_GetCommandStatus {
	c := _m.On("GetCommandStatus", ctx, commandID)
	return &CommandClient_GetCommandStatus{Call: c}
}

func (_m *CommandClient) OnGetCommandStatusMatch(matchers ...interface{}) *CommandClient_GetCommandStatus {
	c := _m.On("GetCommandStatus", matchers...)
	return &CommandClient_GetCommandStatus{Call: c}
}

// GetCommandStatus provides a mock function with given fields: ctx, commandID
func (_m *CommandClient) GetCommandStatus(ctx context.Context, commandID string) (command.CommandStatus, error) {
	ret := _m.Called(ctx, commandID)

	var r0 command.CommandStatus
	if rf, ok := ret.Get(0).(func(context.Context, string) command.CommandStatus); ok {
		r0 = rf(ctx, commandID)
	} else {
		r0 = ret.Get(0).(command.CommandStatus)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, commandID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type CommandClient_KillCommand struct {
	*mock.Call
}

func (_m CommandClient_KillCommand) Return(_a0 error) *CommandClient_KillCommand {
	return &CommandClient_KillCommand{Call: _m.Call.Return(_a0)}
}

func (_m *CommandClient) OnKillCommand(ctx context.Context, commandID string) *CommandClient_KillCommand {
	c := _m.On("KillCommand", ctx, commandID)
	return &CommandClient_KillCommand{Call: c}
}

func (_m *CommandClient) OnKillCommandMatch(matchers ...interface{}) *CommandClient_KillCommand {
	c := _m.On("KillCommand", matchers...)
	return &CommandClient_KillCommand{Call: c}
}

// KillCommand provides a mock function with given fields: ctx, commandID
func (_m *CommandClient) KillCommand(ctx context.Context, commandID string) error {
	ret := _m.Called(ctx, commandID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, commandID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}