package sagemaker

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	commonv1 "github.com/aws/amazon-sagemaker-operator-for-k8s/api/v1/common"
	"github.com/lyft/flyteplugins/go/tasks/plugins/k8s/sagemaker/config"

	sagemakerIdl "github.com/lyft/flyteidl/gen/pb-go/flyteidl/plugins/sagemaker"
	"github.com/stretchr/testify/assert"

	trainingjobv1 "github.com/aws/amazon-sagemaker-operator-for-k8s/api/v1/trainingjob"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	pluginsCore "github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	stdConfig "github.com/lyft/flytestdlib/config"
	"github.com/lyft/flytestdlib/config/viper"
)

func Test_awsSagemakerPlugin_BuildResourceForCustomTrainingJob(t *testing.T) {
	// Default config does not contain a roleAnnotationKey -> expecting to get the role from default config
	ctx := context.TODO()
	defaultCfg := config.GetSagemakerConfig()
	defer func() {
		_ = config.SetSagemakerConfig(defaultCfg)
	}()
	t.Run("In a custom training job we should see the FLYTE_SAGEMAKER_CMD being injected", func(t *testing.T) {
		// Injecting a config which contains a mismatched roleAnnotationKey -> expecting to get the role from the config
		configAccessor := viper.NewAccessor(stdConfig.Options{
			StrictMode: true,
			// Use a different
			SearchPaths: []string{"testdata/config2.yaml"},
		})

		err := configAccessor.UpdateConfig(context.TODO())
		assert.NoError(t, err)

		awsSageMakerTrainingJobHandler := awsSagemakerPlugin{TaskType: customTrainingJobTaskType}

		tjObj := generateMockTrainingJobCustomObj(
			sagemakerIdl.InputMode_FILE, sagemakerIdl.AlgorithmName_CUSTOM, "0.90", []*sagemakerIdl.MetricDefinition{},
			sagemakerIdl.InputContentType_TEXT_CSV, 1, "ml.m4.xlarge", 25)
		taskTemplate := generateMockTrainingJobTaskTemplate("the job", tjObj)

		trainingJobResource, err := awsSageMakerTrainingJobHandler.BuildResource(ctx, generateMockCustomTrainingJobTaskContext(taskTemplate, false))
		assert.NoError(t, err)
		assert.NotNil(t, trainingJobResource)

		trainingJob, ok := trainingJobResource.(*trainingjobv1.TrainingJob)
		assert.True(t, ok)
		assert.Equal(t, "config_role", *trainingJob.Spec.RoleArn)
		//assert.Equal(t, 1, len(trainingJob.Spec.HyperParameters))
		fmt.Printf("%v", trainingJob.Spec.HyperParameters)
		expectedHPs := []*commonv1.KeyValuePair{
			{Name: fmt.Sprintf("%s%d_%s%s", FlyteSageMakerCmdKeyPrefix, 0, "service_venv", FlyteSageMakerKeySuffix), Value: FlyteSageMakerCmdDummyValue},
			{Name: fmt.Sprintf("%s%d_%s%s", FlyteSageMakerCmdKeyPrefix, 1, "pyflyte-execute", FlyteSageMakerKeySuffix), Value: FlyteSageMakerCmdDummyValue},
			{Name: fmt.Sprintf("%s%d_%s%s", FlyteSageMakerCmdKeyPrefix, 2, "--test-opt1", FlyteSageMakerKeySuffix), Value: FlyteSageMakerCmdDummyValue},
			{Name: fmt.Sprintf("%s%d_%s%s", FlyteSageMakerCmdKeyPrefix, 3, "value1", FlyteSageMakerKeySuffix), Value: FlyteSageMakerCmdDummyValue},
			{Name: fmt.Sprintf("%s%d_%s%s", FlyteSageMakerCmdKeyPrefix, 4, "--test-opt2", FlyteSageMakerKeySuffix), Value: FlyteSageMakerCmdDummyValue},
			{Name: fmt.Sprintf("%s%d_%s%s", FlyteSageMakerCmdKeyPrefix, 5, "value2", FlyteSageMakerKeySuffix), Value: FlyteSageMakerCmdDummyValue},
			{Name: fmt.Sprintf("%s%d_%s%s", FlyteSageMakerCmdKeyPrefix, 6, "--test-flag", FlyteSageMakerKeySuffix), Value: FlyteSageMakerCmdDummyValue},
			{Name: fmt.Sprintf("%s%s%s", FlyteSageMakerEnvVarKeyPrefix, "Env_Var", FlyteSageMakerKeySuffix), Value: "Env_Val"},
			{Name: fmt.Sprintf("%s%s%s", FlyteSageMakerEnvVarKeyPrefix, FlyteSageMakerEnvVarKeyStatsdDisabled, FlyteSageMakerKeySuffix), Value: strconv.FormatBool(true)},
		}
		assert.Equal(t, len(expectedHPs), len(trainingJob.Spec.HyperParameters))
		for i := range expectedHPs {
			assert.Equal(t, expectedHPs[i].Name, trainingJob.Spec.HyperParameters[i].Name)
			assert.Equal(t, expectedHPs[i].Value, trainingJob.Spec.HyperParameters[i].Value)
		}

		assert.Equal(t, testImage, *trainingJob.Spec.AlgorithmSpecification.TrainingImage)
	})
}

func Test_awsSagemakerPlugin_GetTaskPhaseForCustomTrainingJob(t *testing.T) {
	ctx := context.TODO()
	// Injecting a config which contains a mismatched roleAnnotationKey -> expecting to get the role from the config
	configAccessor := viper.NewAccessor(stdConfig.Options{
		StrictMode: true,
		// Use a different
		SearchPaths: []string{"testdata/config2.yaml"},
	})

	err := configAccessor.UpdateConfig(context.TODO())
	assert.NoError(t, err)

	awsSageMakerTrainingJobHandler := awsSagemakerPlugin{TaskType: customTrainingJobTaskType}

	t.Run("TrainingJobStatusCompleted", func(t *testing.T) {
		tjObj := generateMockTrainingJobCustomObj(
			sagemakerIdl.InputMode_FILE, sagemakerIdl.AlgorithmName_XGBOOST, "0.90", []*sagemakerIdl.MetricDefinition{},
			sagemakerIdl.InputContentType_TEXT_CSV, 1, "ml.m4.xlarge", 25)
		taskTemplate := generateMockTrainingJobTaskTemplate("the job", tjObj)
		taskCtx := generateMockCustomTrainingJobTaskContext(taskTemplate, false)
		trainingJobResource, err := awsSageMakerTrainingJobHandler.BuildResource(ctx, taskCtx)
		assert.Error(t, err)
		assert.Nil(t, trainingJobResource)
	})

	t.Run("TrainingJobStatusCompleted", func(t *testing.T) {
		tjObj := generateMockTrainingJobCustomObj(
			sagemakerIdl.InputMode_FILE, sagemakerIdl.AlgorithmName_CUSTOM, "", []*sagemakerIdl.MetricDefinition{},
			sagemakerIdl.InputContentType_TEXT_CSV, 1, "ml.m4.xlarge", 25)
		taskTemplate := generateMockTrainingJobTaskTemplate("the job", tjObj)
		taskCtx := generateMockCustomTrainingJobTaskContext(taskTemplate, false)
		trainingJobResource, err := awsSageMakerTrainingJobHandler.BuildResource(ctx, taskCtx)
		assert.NoError(t, err)
		assert.NotNil(t, trainingJobResource)

		trainingJob, ok := trainingJobResource.(*trainingjobv1.TrainingJob)
		assert.True(t, ok)

		trainingJob.Status.TrainingJobStatus = sagemaker.TrainingJobStatusCompleted
		phaseInfo, err := awsSageMakerTrainingJobHandler.getTaskPhaseForCustomTrainingJob(ctx, taskCtx, trainingJob)
		assert.Nil(t, err)
		assert.Equal(t, phaseInfo.Phase(), pluginsCore.PhaseSuccess)
	})
	t.Run("OutputWriter.Put returns an error", func(t *testing.T) {
		tjObj := generateMockTrainingJobCustomObj(
			sagemakerIdl.InputMode_FILE, sagemakerIdl.AlgorithmName_CUSTOM, "", []*sagemakerIdl.MetricDefinition{},
			sagemakerIdl.InputContentType_TEXT_CSV, 1, "ml.m4.xlarge", 25)
		taskTemplate := generateMockTrainingJobTaskTemplate("the job", tjObj)
		taskCtx := generateMockCustomTrainingJobTaskContext(taskTemplate, true)
		trainingJobResource, err := awsSageMakerTrainingJobHandler.BuildResource(ctx, taskCtx)
		assert.NoError(t, err)
		assert.NotNil(t, trainingJobResource)

		trainingJob, ok := trainingJobResource.(*trainingjobv1.TrainingJob)
		assert.True(t, ok)

		trainingJob.Status.TrainingJobStatus = sagemaker.TrainingJobStatusCompleted
		phaseInfo, err := awsSageMakerTrainingJobHandler.getTaskPhaseForCustomTrainingJob(ctx, taskCtx, trainingJob)
		assert.NotNil(t, err)
		assert.Equal(t, phaseInfo.Phase(), pluginsCore.PhaseUndefined)
	})
}
