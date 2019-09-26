package hive

import (
	"context"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core/mocks"
	"github.com/lyft/flyteplugins/go/tasks/plugins/hive/client"
	quboleMocks "github.com/lyft/flyteplugins/go/tasks/plugins/hive/client/mocks"
	stdlibMocks "github.com/lyft/flytestdlib/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"testing"
)

func TestInTerminalState(t *testing.T) {
	var stateTests = []struct {
		phase      ExecutionPhase
		isTerminal bool
	}{
		{phase: PhaseNotStarted, isTerminal: false},
		{phase: PhaseQueued, isTerminal: false},
		{phase: PhaseSubmitted, isTerminal: false},
		{phase: PhaseQuerySucceeded, isTerminal: true},
		{phase: PhaseQueryFailed, isTerminal: true},
	}

	for _, tt := range stateTests {
		t.Run(tt.phase.String(), func(t *testing.T) {
			e := ExecutionState{Phase: tt.phase}
			res := InTerminalState(e)
			assert.Equal(t, tt.isTerminal, res)
		})
	}
}

func TestIsNotYetSubmitted(t *testing.T) {
	var stateTests = []struct {
		phase             ExecutionPhase
		isNotYetSubmitted bool
	}{
		{phase: PhaseNotStarted, isNotYetSubmitted: true},
		{phase: PhaseQueued, isNotYetSubmitted: true},
		{phase: PhaseSubmitted, isNotYetSubmitted: false},
		{phase: PhaseQuerySucceeded, isNotYetSubmitted: false},
		{phase: PhaseQueryFailed, isNotYetSubmitted: false},
	}

	for _, tt := range stateTests {
		t.Run(tt.phase.String(), func(t *testing.T) {
			e := ExecutionState{Phase: tt.phase}
			res := IsNotYetSubmitted(e)
			assert.Equal(t, tt.isNotYetSubmitted, res)
		})
	}
}

func TestGetQueryInfo(t *testing.T) {
	ctx := context.Background()

	taskTemplate := GetSingleHiveQueryTaskTemplate()
	mockTaskReader := &mocks.TaskReader{}
	mockTaskReader.On("Read", mock.Anything).Return(&taskTemplate, nil)

	mockTaskExecutionContext := mocks.TaskExecutionContext{}
	mockTaskExecutionContext.On("TaskReader").Return(mockTaskReader)

	query, cluster, tags, timeout, err := GetQueryInfo(ctx, &mockTaskExecutionContext)
	assert.NoError(t, err)
	assert.Equal(t, "select 'one'", query)
	assert.Equal(t, "default", cluster)
	assert.Equal(t, []string{"flyte_plugin_test"}, tags)
	assert.Equal(t, 500, int(timeout))
}

func TestConstructTaskLog(t *testing.T) {
	taskLog := ConstructTaskLog(ExecutionState{CommandId: "123"})
	assert.Equal(t, "https://api.qubole.com/v2/analyze?command_id=123", taskLog.Uri)
}

func TestConstructTaskInfo(t *testing.T) {
	empty := ConstructTaskInfo(ExecutionState{})
	assert.Nil(t, empty)

	e := ExecutionState{
		Phase:            PhaseQuerySucceeded,
		CommandId:        "123",
		SyncFailureCount: 0,
	}
	taskInfo := ConstructTaskInfo(e)
	assert.Equal(t, "https://api.qubole.com/v2/analyze?command_id=123", taskInfo.Logs[0].Uri)
}

func TestMapExecutionStateToPhaseInfo(t *testing.T) {
	t.Run("NotStarted", func(t *testing.T) {
		e := ExecutionState{
			Phase: PhaseNotStarted,
		}
		phaseInfo := MapExecutionStateToPhaseInfo(e)
		assert.Equal(t, core.PhaseNotReady, phaseInfo.Phase())
	})

	t.Run("Queued", func(t *testing.T) {
		e := ExecutionState{
			Phase:                PhaseQueued,
			CreationFailureCount: 0,
		}
		phaseInfo := MapExecutionStateToPhaseInfo(e)
		assert.Equal(t, core.PhaseQueued, phaseInfo.Phase())

		e = ExecutionState{
			Phase:                PhaseQueued,
			CreationFailureCount: 100,
		}
		phaseInfo = MapExecutionStateToPhaseInfo(e)
		assert.Equal(t, core.PhaseRetryableFailure, phaseInfo.Phase())

	})

	t.Run("Submitted", func(t *testing.T) {
		e := ExecutionState{
			Phase: PhaseSubmitted,
		}
		phaseInfo := MapExecutionStateToPhaseInfo(e)
		assert.Equal(t, core.PhaseRunning, phaseInfo.Phase())
	})
}

func TestGetAllocationToken(t *testing.T) {
	ctx := context.Background()

	t.Run("allocation granted", func(t *testing.T) {
		tCtx := GetMockTaskExecutionContext()
		mockResourceManager := tCtx.ResourceManager()
		x := mockResourceManager.(*mocks.ResourceManager)
		x.On("AllocateResource", mock.Anything, mock.Anything, mock.Anything).
			Return(core.AllocationStatusGranted, nil)

		state, err := GetAllocationToken(ctx, tCtx)
		assert.NoError(t, err)
		assert.Equal(t, PhaseQueued, state.Phase)
	})

	t.Run("exhausted", func(t *testing.T) {
		tCtx := GetMockTaskExecutionContext()
		mockResourceManager := tCtx.ResourceManager()
		x := mockResourceManager.(*mocks.ResourceManager)
		x.On("AllocateResource", mock.Anything, mock.Anything, mock.Anything).
			Return(core.AllocationStatusExhausted, nil)

		state, err := GetAllocationToken(ctx, tCtx)
		assert.NoError(t, err)
		assert.Equal(t, PhaseNotStarted, state.Phase)
	})

	t.Run("namespace exhausted", func(t *testing.T) {
		tCtx := GetMockTaskExecutionContext()
		mockResourceManager := tCtx.ResourceManager()
		x := mockResourceManager.(*mocks.ResourceManager)
		x.On("AllocateResource", mock.Anything, mock.Anything, mock.Anything).
			Return(core.AllocationStatusNamespaceQuotaExceeded, nil)

		state, err := GetAllocationToken(ctx, tCtx)
		assert.NoError(t, err)
		assert.Equal(t, PhaseNotStarted, state.Phase)
	})
}

func TestAbort(t *testing.T) {
	ctx := context.Background()

	t.Run("Terminate called when not in terminal state", func(t *testing.T) {
		var x = false
		mockQubole := &quboleMocks.QuboleClient{}
		mockQubole.On("KillCommand", mock.Anything, mock.MatchedBy(func(commandId string) bool {
			return commandId == "123456"
		}), mock.Anything).Run(func(_ mock.Arguments) {
			x = true
		}).Return(nil)

		err := Abort(ctx, GetMockTaskExecutionContext(), ExecutionState{Phase: PhaseSubmitted, CommandId: "123456"}, mockQubole)
		assert.NoError(t, err)
		assert.True(t, x)
	})

	t.Run("Terminate not called when in terminal state", func(t *testing.T) {
		var x = false
		mockQubole := &quboleMocks.QuboleClient{}
		mockQubole.On("KillCommand", mock.Anything, mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
			x = true
		}).Return(nil)

		err := Abort(ctx, GetMockTaskExecutionContext(), ExecutionState{
			Phase:     PhaseQuerySucceeded,
			CommandId: "123456",
		}, mockQubole)
		assert.NoError(t, err)
		assert.False(t, x)
	})
}

func TestFinalize(t *testing.T) {
	// Test that Finalize releases resources
	ctx := context.Background()
	tCtx := GetMockTaskExecutionContext()
	state := ExecutionState{}
	var called = false
	mockResourceManager := tCtx.ResourceManager()
	x := mockResourceManager.(*mocks.ResourceManager)
	x.On("ReleaseResource", mock.Anything, mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		called = true
	}).Return(nil)

	err := Finalize(ctx, tCtx, state)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestMonitorQuery(t *testing.T) {
	ctx := context.Background()
	tCtx := GetMockTaskExecutionContext()
	state := ExecutionState{
		Phase: PhaseSubmitted,
	}
	var getOrCreateCalled = false
	mockCache := &stdlibMocks.AutoRefreshCache{}
	mockCache.On("GetOrCreate", mock.Anything).Run(func(_ mock.Arguments) {
		getOrCreateCalled = true
	}).Return(ExecutionStateCacheItem{
		ExecutionState: ExecutionState{Phase: PhaseQuerySucceeded},
		Id:             "some-id",
	}, nil)

	newState, err := MonitorQuery(ctx, tCtx, state, mockCache)
	assert.NoError(t, err)
	assert.True(t, getOrCreateCalled)
	assert.Equal(t, PhaseQuerySucceeded, newState.Phase)
}

func TestKickOffQuery(t *testing.T) {
	ctx := context.Background()
	tCtx := GetMockTaskExecutionContext()

	var quboleCalled = false
	quboleCommandDetails := &client.QuboleCommandDetails{
		ID:     int64(453298043),
		Status: client.QuboleStatusWaiting,
	}
	mockQubole := &quboleMocks.QuboleClient{}
	mockQubole.On("ExecuteHiveCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		quboleCalled = true
	}).Return(quboleCommandDetails, nil)

	var getOrCreateCalled = false
	mockCache := &stdlibMocks.AutoRefreshCache{}
	mockCache.On("GetOrCreate", mock.Anything).Run(func(_ mock.Arguments) {
		getOrCreateCalled = true
	}).Return(ExecutionStateCacheItem{}, nil)

	state := ExecutionState{}
	newState, err := KickOffQuery(ctx, tCtx, state, mockQubole, mockCache)
	assert.NoError(t, err)
	assert.Equal(t, PhaseSubmitted, newState.Phase)
	assert.Equal(t, "453298043", newState.CommandId)
	assert.True(t, getOrCreateCalled)
	assert.True(t, quboleCalled)
}
