// Code generated by mockery v1.0.1. DO NOT EDIT.

package mocks

import (
	core "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	io "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/io"

	mock "github.com/stretchr/testify/mock"
)

// TaskExecutionContextReader is an autogenerated mock type for the TaskExecutionContextReader type
type TaskExecutionContextReader struct {
	mock.Mock
}

type TaskExecutionContextReader_InputReader struct {
	*mock.Call
}

func (_m TaskExecutionContextReader_InputReader) Return(_a0 io.InputReader) *TaskExecutionContextReader_InputReader {
	return &TaskExecutionContextReader_InputReader{Call: _m.Call.Return(_a0)}
}

func (_m *TaskExecutionContextReader) OnInputReader() *TaskExecutionContextReader_InputReader {
	c := _m.On("InputReader")
	return &TaskExecutionContextReader_InputReader{Call: c}
}

func (_m *TaskExecutionContextReader) OnInputReaderMatch(matchers ...interface{}) *TaskExecutionContextReader_InputReader {
	c := _m.On("InputReader", matchers...)
	return &TaskExecutionContextReader_InputReader{Call: c}
}

// InputReader provides a mock function with given fields:
func (_m *TaskExecutionContextReader) InputReader() io.InputReader {
	ret := _m.Called()

	var r0 io.InputReader
	if rf, ok := ret.Get(0).(func() io.InputReader); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.InputReader)
		}
	}

	return r0
}

type TaskExecutionContextReader_TaskExecutionMetadata struct {
	*mock.Call
}

func (_m TaskExecutionContextReader_TaskExecutionMetadata) Return(_a0 core.TaskExecutionMetadata) *TaskExecutionContextReader_TaskExecutionMetadata {
	return &TaskExecutionContextReader_TaskExecutionMetadata{Call: _m.Call.Return(_a0)}
}

func (_m *TaskExecutionContextReader) OnTaskExecutionMetadata() *TaskExecutionContextReader_TaskExecutionMetadata {
	c := _m.On("TaskExecutionMetadata")
	return &TaskExecutionContextReader_TaskExecutionMetadata{Call: c}
}

func (_m *TaskExecutionContextReader) OnTaskExecutionMetadataMatch(matchers ...interface{}) *TaskExecutionContextReader_TaskExecutionMetadata {
	c := _m.On("TaskExecutionMetadata", matchers...)
	return &TaskExecutionContextReader_TaskExecutionMetadata{Call: c}
}

// TaskExecutionMetadata provides a mock function with given fields:
func (_m *TaskExecutionContextReader) TaskExecutionMetadata() core.TaskExecutionMetadata {
	ret := _m.Called()

	var r0 core.TaskExecutionMetadata
	if rf, ok := ret.Get(0).(func() core.TaskExecutionMetadata); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.TaskExecutionMetadata)
		}
	}

	return r0
}

type TaskExecutionContextReader_TaskReader struct {
	*mock.Call
}

func (_m TaskExecutionContextReader_TaskReader) Return(_a0 core.TaskReader) *TaskExecutionContextReader_TaskReader {
	return &TaskExecutionContextReader_TaskReader{Call: _m.Call.Return(_a0)}
}

func (_m *TaskExecutionContextReader) OnTaskReader() *TaskExecutionContextReader_TaskReader {
	c := _m.On("TaskReader")
	return &TaskExecutionContextReader_TaskReader{Call: c}
}

func (_m *TaskExecutionContextReader) OnTaskReaderMatch(matchers ...interface{}) *TaskExecutionContextReader_TaskReader {
	c := _m.On("TaskReader", matchers...)
	return &TaskExecutionContextReader_TaskReader{Call: c}
}

// TaskReader provides a mock function with given fields:
func (_m *TaskExecutionContextReader) TaskReader() core.TaskReader {
	ret := _m.Called()

	var r0 core.TaskReader
	if rf, ok := ret.Get(0).(func() core.TaskReader); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.TaskReader)
		}
	}

	return r0
}
