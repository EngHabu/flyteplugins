package k8s

import (
	"context"
	"github.com/lyft/flyteplugins/go/tasks/plugins/array"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery"

	"github.com/lyft/flyteplugins/go/tasks/errors"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
)

const executorName = "k8s-array"
const arrayTaskType = "container_array"
const pluginStateVersion = 0

type Executor struct {
	kubeClient       core.KubeClient
	outputsAssembler array.OutputAssembler
	errorAssembler   array.OutputAssembler
}

func NewExecutor(kubeClient core.KubeClient, cfg *Config) (Executor, error) {
	outputAssembler, err := array.NewOutputAssembler(cfg.OutputAssembler)
	if err != nil {
		return Executor{}, err
	}

	errorAssembler, err := array.NewErrorAssembler(cfg.MaxErrorStringLength, cfg.ErrorAssembler)
	if err != nil {
		return Executor{}, err
	}

	return Executor{
		kubeClient:       kubeClient,
		outputsAssembler: outputAssembler,
		errorAssembler:   errorAssembler,
	}, nil
}

func (e Executor) GetID() string {
	return executorName
}

func (Executor) GetProperties() core.PluginProperties {
	return core.PluginProperties{}
}

func (e Executor) Handle(ctx context.Context, tCtx core.TaskExecutionContext) (core.Transition, error) {
	pluginConfig := GetConfig()

	pluginState := &array.StateImpl{}
	if _, err := tCtx.PluginStateReader().Get(pluginState); err != nil {
		return core.UnknownTransition, errors.Wrapf(errors.CorruptedPluginState, err, "Failed to read unmarshal custom state")
	}

	var nextState array.State
	var err error

	switch p, _ := pluginState.GetPhase(); p {
	case array.PhaseStart:
		nextState, err = array.DetermineDiscoverability(ctx, tCtx, pluginState)

	case array.PhasePreLaunch:
		nextState = pluginState.SetPhase(array.PhaseLaunch, core.DefaultPhaseVersion)
		err = nil

	case array.PhaseLaunch:
		nextState, err = LaunchSubTasks(ctx, tCtx, e.kubeClient, pluginConfig, pluginState)

	case array.PhaseCheckingSubTaskExecutions:
		nextState, err = CheckSubTasksState(ctx, tCtx, e.kubeClient, pluginConfig, pluginState)

	case array.PhaseAssembleFinalOutput:
		nextState, err = array.AssembleFinalOutputs(ctx, e.outputsAssembler, tCtx, pluginState)

	case array.PhaseWriteToDiscovery:
		nextState, err = array.WriteToDiscovery(ctx, tCtx, pluginState)

	case array.PhaseAssembleFinalError:
		nextState, err = array.AssembleFinalOutputs(ctx, e.errorAssembler, tCtx, pluginState)

	default:
		nextState = pluginState
		err = nil
	}
	if err != nil {
		return core.UnknownTransition, err
	}

	if err := tCtx.PluginStateWriter().Put(pluginStateVersion, nextState); err != nil {
		return core.UnknownTransition, err
	}

	// Determine transition information from the state
	phaseInfo := array.MapArrayStateToPluginPhase(ctx, nextState)
	return core.DoTransitionType(core.TransitionTypeBestEffort, phaseInfo), nil
}

func (Executor) Abort(ctx context.Context, tCtx core.TaskExecutionContext) error {
	return nil
}

func (Executor) Finalize(ctx context.Context, tCtx core.TaskExecutionContext) error {
	return nil
}

func init() {
	pluginmachinery.PluginRegistry().RegisterCorePlugin(
		core.PluginEntry{
			ID:                  executorName,
			RegisteredTaskTypes: []core.TaskType{arrayTaskType},
			LoadPlugin:          GetNewExecutorPlugin,
			IsDefault:           false,
		})
}

func GetNewExecutorPlugin(_ context.Context, iCtx core.SetupContext) (core.Plugin, error) {
	return NewExecutor(iCtx.KubeClient(), GetConfig())
}