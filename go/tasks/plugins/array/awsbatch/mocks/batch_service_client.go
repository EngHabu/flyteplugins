// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import batch "github.com/aws/aws-sdk-go/service/batch"
import context "context"
import mock "github.com/stretchr/testify/mock"
import request "github.com/aws/aws-sdk-go/aws/request"

// BatchServiceClient is an autogenerated mock type for the BatchServiceClient type
type BatchServiceClient struct {
	mock.Mock
}

// DescribeJobsWithContext provides a mock function with given fields: ctx, input, opts
func (_m *BatchServiceClient) DescribeJobsWithContext(ctx context.Context, input *batch.DescribeJobsInput, opts ...request.Option) (*batch.DescribeJobsOutput, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, input)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *batch.DescribeJobsOutput
	if rf, ok := ret.Get(0).(func(context.Context, *batch.DescribeJobsInput, ...request.Option) *batch.DescribeJobsOutput); ok {
		r0 = rf(ctx, input, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*batch.DescribeJobsOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *batch.DescribeJobsInput, ...request.Option) error); ok {
		r1 = rf(ctx, input, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterJobDefinitionWithContext provides a mock function with given fields: ctx, input, opts
func (_m *BatchServiceClient) RegisterJobDefinitionWithContext(ctx context.Context, input *batch.RegisterJobDefinitionInput, opts ...request.Option) (*batch.RegisterJobDefinitionOutput, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, input)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *batch.RegisterJobDefinitionOutput
	if rf, ok := ret.Get(0).(func(context.Context, *batch.RegisterJobDefinitionInput, ...request.Option) *batch.RegisterJobDefinitionOutput); ok {
		r0 = rf(ctx, input, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*batch.RegisterJobDefinitionOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *batch.RegisterJobDefinitionInput, ...request.Option) error); ok {
		r1 = rf(ctx, input, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubmitJobWithContext provides a mock function with given fields: ctx, input, opts
func (_m *BatchServiceClient) SubmitJobWithContext(ctx context.Context, input *batch.SubmitJobInput, opts ...request.Option) (*batch.SubmitJobOutput, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, input)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *batch.SubmitJobOutput
	if rf, ok := ret.Get(0).(func(context.Context, *batch.SubmitJobInput, ...request.Option) *batch.SubmitJobOutput); ok {
		r0 = rf(ctx, input, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*batch.SubmitJobOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *batch.SubmitJobInput, ...request.Option) error); ok {
		r1 = rf(ctx, input, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TerminateJobWithContext provides a mock function with given fields: ctx, input, opts
func (_m *BatchServiceClient) TerminateJobWithContext(ctx context.Context, input *batch.TerminateJobInput, opts ...request.Option) (*batch.TerminateJobOutput, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, input)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *batch.TerminateJobOutput
	if rf, ok := ret.Get(0).(func(context.Context, *batch.TerminateJobInput, ...request.Option) *batch.TerminateJobOutput); ok {
		r0 = rf(ctx, input, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*batch.TerminateJobOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *batch.TerminateJobInput, ...request.Option) error); ok {
		r1 = rf(ctx, input, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
