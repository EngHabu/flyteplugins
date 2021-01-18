package logs

import (
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/lyft/flyteplugins/go/tasks/config"
)

//go:generate pflags LogConfig

// A URI that accepts templates. See: go/tasks/pluginmachinery/tasklog/template.go for available templates.
type TemplateURI = string

// Log plugins configs
type LogConfig struct {
	IsCloudwatchEnabled bool `json:"cloudwatch-enabled" pflag:",Enable Cloudwatch Logging"`
	// Deprecated: Please use CloudwatchTemplateUri
	CloudwatchRegion string `json:"cloudwatch-region" pflag:",AWS region in which Cloudwatch logs are stored."`
	// Deprecated: Please use CloudwatchTemplateUri
	CloudwatchLogGroup    string      `json:"cloudwatch-log-group" pflag:",Log group to which streams are associated."`
	CloudwatchTemplateUri TemplateURI `json:"cloudwatch-template-uri" pflag:",Template Uri to use when building cloudwatch log links"`

	IsKubernetesEnabled bool `json:"kubernetes-enabled" pflag:",Enable Kubernetes Logging"`
	// Deprecated: Please use KubernetesTemplateUri
	KubernetesURL         string      `json:"kubernetes-url" pflag:",Console URL for Kubernetes logs"`
	KubernetesTemplateUri TemplateURI `json:"kubernetes-template-uri" pflag:",Template Uri to use when building kubernetes log links"`

	IsStackDriverEnabled bool `json:"stackdriver-enabled" pflag:",Enable Log-links to stackdriver"`
	// Deprecated: Please use StackDriverTemplateUri
	GCPProjectName string `json:"gcp-project" pflag:",Name of the project in GCP"`
	// Deprecated: Please use StackDriverTemplateUri
	StackdriverLogResourceName string      `json:"stackdriver-logresourcename" pflag:",Name of the logresource in stackdriver"`
	StackDriverTemplateUri     TemplateURI `json:"stackdriver-template-uri" pflag:",Template Uri to use when building stackdriver log links"`

	Templates []TemplateLogPluginConfig `json:"templates" pflag:","`
}

type TemplateLogPluginConfig struct {
	DisplayName   string                     `json:"displayName" pflag:",Display name for the generated log when displayed in the console."`
	TemplateURIs  []TemplateURI              `json:"templateUris" pflag:",URI Templates for generating task log links."`
	MessageFormat core.TaskLog_MessageFormat `json:"messageFormat" pflag:",Log Message Format."`
}

var (
	logConfigSection = config.MustRegisterSubSection("logs", &LogConfig{})
)

func GetLogConfig() *LogConfig {
	return logConfigSection.GetConfig().(*LogConfig)
}

// This method should be used for unit testing only
func SetLogConfig(logConfig *LogConfig) error {
	return logConfigSection.SetConfig(logConfig)
}
