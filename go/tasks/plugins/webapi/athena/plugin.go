package athena

import (
	"context"
	"fmt"
	"time"

	pluginsIdl "github.com/lyft/flyteidl/gen/pb-go/flyteidl/plugins"

	awsSdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	athenaTypes "github.com/aws/aws-sdk-go-v2/service/athena/types"
	"github.com/lyft/flyteplugins/go/tasks/aws"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/ioutils"

	idlCore "github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/lyft/flytestdlib/errors"
	"github.com/lyft/flytestdlib/utils"

	"github.com/lyft/flytestdlib/logger"

	"github.com/lyft/flytestdlib/promutils"

	pb "github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/webapi"
)

const (
	ErrRemoteSystem errors.ErrorCode = "RemoteSystem"
	ErrRemoteUser   errors.ErrorCode = "RemoteUser"
	ErrSystem       errors.ErrorCode = "System"
)

type Plugin struct {
	metricScope promutils.Scope
	client      *athena.Client
	cfg         *Config
	awsConfig   awsSdk.Config
}

type ResourceWrapper struct {
	Status               *athenaTypes.QueryExecutionStatus
	ResultsConfiguration *athenaTypes.ResultConfiguration
}

func (p Plugin) GetConfig() webapi.PluginConfig {
	return GetConfig().WebAPI
}

func (p Plugin) ResourceRequirements(_ context.Context, _ webapi.TaskExecutionContextReader) (
	namespace core.ResourceNamespace, constraints core.ResourceConstraintsSpec, err error) {

	// Resource requirements are assumed to be the same.
	return "default", p.cfg.ResourceConstraints, nil
}

func (p Plugin) Create(ctx context.Context, tCtx webapi.TaskExecutionContextReader) (resourceMeta webapi.ResourceMeta,
	resource webapi.Resource, err error) {

	// TODO: explain what this block does... and why...
	// TODO: open an issue to add ReadCustom()
	task, err := tCtx.TaskReader().Read(ctx)
	if err != nil {
		return nil, nil, err
	}

	custom := task.GetCustom()
	hiveQuery := &pluginsIdl.QuboleHiveJob{}
	err = utils.UnmarshalStructToPb(custom, hiveQuery)
	if err != nil {
		return nil, nil, err
	}

	if hiveQuery.Query == nil {
		return "", "", nil
	}

	execID := tCtx.TaskExecutionMetadata().GetTaskExecutionID().GetID().NodeExecutionId.GetExecutionId()
	resp, err := p.client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
		ClientRequestToken: awsSdk.String(fmt.Sprintf("%v-%v-%v", execID.Project, execID.Domain, tCtx.TaskExecutionMetadata().GetTaskExecutionID().GetGeneratedName())),
		QueryExecutionContext: &athenaTypes.QueryExecutionContext{
			Database: awsSdk.String("vaccinations"),
			Catalog:  awsSdk.String(p.cfg.DefaultCatalog),
			//Description: awsSdk.String(fmt.Sprintf("Launched query through Athena Plugin for execution [%v]", execID)),
			//Name:        awsSdk.String(tCtx.TaskExecutionMetadata().GetTaskExecutionID().GetGeneratedName()),
		},
		ResultConfiguration: &athenaTypes.ResultConfiguration{
			// Workgroup settings can override the output location setting.
			OutputLocation: awsSdk.String(tCtx.OutputWriter().GetRawOutputPrefix().String()),
		},
		QueryString: awsSdk.String(hiveQuery.Query.Query),
		WorkGroup:   awsSdk.String(p.cfg.DefaultWorkGroup),
	})

	if err != nil {
		return "", "", err
	}

	if resp.QueryExecutionId == nil {
		return "", "", errors.Errorf(ErrRemoteSystem, "Service created an empty query id")
	}

	return *resp.QueryExecutionId, nil, nil
}

func (p Plugin) Get(ctx context.Context, tCtx webapi.GetContext) (latest webapi.Resource, err error) {
	exec := tCtx.ResourceMeta().(string)
	resp, err := p.client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{
		QueryExecutionId: awsSdk.String(exec),
	})
	if err != nil {
		return nil, err
	}

	// Only cache fields we want to keep in memory instead of the potentially huge execution closure.
	return ResourceWrapper{
		Status:               resp.QueryExecution.Status,
		ResultsConfiguration: resp.QueryExecution.ResultConfiguration,
	}, nil
}

func (p Plugin) Delete(ctx context.Context, tCtx webapi.DeleteContext) error {
	resp, err := p.client.StopQueryExecution(ctx, &athena.StopQueryExecutionInput{
		QueryExecutionId: awsSdk.String(tCtx.ResourceMeta().(string)),
	})
	if err != nil {
		return err
	}

	logger.Info(ctx, "Deleted query execution [%v]", resp)

	return nil
}

func writeOutput(ctx context.Context, tCtx webapi.StatusContext, externalLocation string) error {
	taskTemplate, err := tCtx.TaskReader().Read(ctx)
	if err != nil {
		return err
	}

	resultsSchema, exists := taskTemplate.Interface.Outputs.Variables["results"]
	if !exists {
		logger.Infof(ctx, "The task declares no outputs. Skipping writing the outputs.")
		return nil
	}

	return tCtx.OutputWriter().Put(ctx, ioutils.NewInMemoryOutputReader(
		&pb.LiteralMap{
			Literals: map[string]*pb.Literal{
				"results": {
					Value: &pb.Literal_Scalar{
						Scalar: &pb.Scalar{Value: &pb.Scalar_Schema{
							Schema: &pb.Schema{
								Uri:  externalLocation,
								Type: resultsSchema.GetType().GetSchema(),
							},
						},
						},
					},
				},
			},
		}, nil))
}

func (p Plugin) Status(ctx context.Context, tCtx webapi.StatusContext) (phase core.PhaseInfo, err error) {
	execID := tCtx.ResourceMeta().(string)
	exec := tCtx.Resource().(ResourceWrapper)
	if exec.Status == nil {
		return core.PhaseInfoUndefined, errors.Errorf(ErrSystem, "No Status field set.")
	}

	switch exec.Status.State {
	case athenaTypes.QueryExecutionStateQueued:
		fallthrough
	case athenaTypes.QueryExecutionStateRunning:
		return core.PhaseInfoRunning(0, createTaskInfo(execID, p.awsConfig)), nil
	case athenaTypes.QueryExecutionStateCancelled:
		reason := "Remote execution was aborted."
		if reasonPtr := exec.Status.StateChangeReason; reasonPtr != nil {
			reason = *reasonPtr
		}

		return core.PhaseInfoRetryableFailure("ABORTED", reason, createTaskInfo(execID, p.awsConfig)), nil
	case athenaTypes.QueryExecutionStateFailed:
		reason := "Remote execution failed"
		if reasonPtr := exec.Status.StateChangeReason; reasonPtr != nil {
			reason = *reasonPtr
		}

		return core.PhaseInfoRetryableFailure("FAILED", reason, createTaskInfo(execID, p.awsConfig)), nil
	case athenaTypes.QueryExecutionStateSucceeded:
		if outputLocation := exec.ResultsConfiguration.OutputLocation; outputLocation != nil {
			// If WorkGroup settings overrode the client settings, the location submitted in the request might have been
			// ignored.
			err = writeOutput(ctx, tCtx, *outputLocation)
			if err != nil {
				logger.Warnf(ctx, "Failed to write output, uri [%s], err %s", *outputLocation, err.Error())
				return core.PhaseInfoUndefined, err
			}
		}

		return core.PhaseInfoSuccess(createTaskInfo(execID, p.awsConfig)), nil
	}

	return core.PhaseInfoUndefined, errors.Errorf(ErrSystem, "Unknown execution phase [%v].", exec.Status.State)
}

func createTaskInfo(queryID string, cfg awsSdk.Config) *core.TaskInfo {
	timeNow := time.Now()
	return &core.TaskInfo{
		OccurredAt: &timeNow,
		Logs: []*idlCore.TaskLog{
			{
				Uri: fmt.Sprintf("https://%v.console.aws.amazon.com/athena/home?force&region=%v#query/history/%v",
					cfg.Region,
					cfg.Region,
					queryID),
				Name: "Athena Query History",
			},
		},
	}
}

func NewPlugin(_ context.Context, cfg *Config, awsConfig *aws.Config, metricScope promutils.Scope) (Plugin, error) {
	sdkCfg, err := awsConfig.GetSdkConfig()
	if err != nil {
		return Plugin{}, err
	}

	return Plugin{
		metricScope: metricScope,
		client:      athena.NewFromConfig(sdkCfg),
		cfg:         cfg,
		awsConfig:   sdkCfg,
	}, nil
}

func init() {
	pluginmachinery.PluginRegistry().RegisterRemotePlugin(webapi.PluginEntry{
		ID:                 "athena",
		SupportedTaskTypes: []core.TaskType{"hive"},
		PluginLoader: func(ctx context.Context, iCtx webapi.PluginSetupContext) (webapi.AsyncPlugin, error) {
			return NewPlugin(ctx, GetConfig(), aws.GetConfig(), iCtx.MetricsScope())
		},
	})
}
