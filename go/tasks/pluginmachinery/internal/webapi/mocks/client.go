// Code generated by mockery v1.0.1. DO NOT EDIT.

package mocks

import (
	context "context"

	core "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

type Client_Get struct {
	*mock.Call
}

func (_m Client_Get) Return(latest interface{}, err error) *Client_Get {
	return &Client_Get{Call: _m.Call.Return(latest, err)}
}

func (_m *Client) OnGet(ctx context.Context, cached interface{}) *Client_Get {
	c := _m.On("Get", ctx, cached)
	return &Client_Get{Call: c}
}

func (_m *Client) OnGetMatch(matchers ...interface{}) *Client_Get {
	c := _m.On("Get", matchers...)
	return &Client_Get{Call: c}
}

// Get provides a mock function with given fields: ctx, cached
func (_m *Client) Get(ctx context.Context, cached interface{}) (interface{}, error) {
	ret := _m.Called(ctx, cached)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) interface{}); ok {
		r0 = rf(ctx, cached)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(ctx, cached)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type Client_Status struct {
	*mock.Call
}

func (_m Client_Status) Return(phase core.PhaseInfo, err error) *Client_Status {
	return &Client_Status{Call: _m.Call.Return(phase, err)}
}

func (_m *Client) OnStatus(ctx context.Context, resource interface{}) *Client_Status {
	c := _m.On("Status", ctx, resource)
	return &Client_Status{Call: c}
}

func (_m *Client) OnStatusMatch(matchers ...interface{}) *Client_Status {
	c := _m.On("Status", matchers...)
	return &Client_Status{Call: c}
}

// Status provides a mock function with given fields: ctx, resource
func (_m *Client) Status(ctx context.Context, resource interface{}) (core.PhaseInfo, error) {
	ret := _m.Called(ctx, resource)

	var r0 core.PhaseInfo
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) core.PhaseInfo); ok {
		r0 = rf(ctx, resource)
	} else {
		r0 = ret.Get(0).(core.PhaseInfo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(ctx, resource)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
