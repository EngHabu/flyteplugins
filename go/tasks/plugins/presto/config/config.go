package config

//go:generate pflags Config --default-var=defaultConfig

import (
	"context"
	"net/url"
	"time"

	"github.com/lyft/flytestdlib/config"
	"github.com/lyft/flytestdlib/logger"

	pluginsConfig "github.com/lyft/flyteplugins/go/tasks/config"
)

const prestoConfigSectionKey = "presto"

func URLMustParse(s string) config.URL {
	r, err := url.Parse(s)
	if err != nil {
		logger.Panicf(context.TODO(), "Bad Presto URL Specified as default, error: %s", err)
	}
	if r == nil {
		logger.Panicf(context.TODO(), "Nil Presto URL specified.", err)
	}
	return config.URL{URL: *r}
}

type RoutingGroupConfig struct {
	Name                             string  `json:"name" pflag:",The name of a given Presto routing group"`
	Limit                            int     `json:"limit" pflag:",Resource quota (in the number of outstanding requests) of the routing group"`
	ProjectScopeQuotaProportionCap   float64 `json:"projectScopeQuotaProportionCap" pflag:",A floating point number between 0 and 1, specifying the maximum proportion of quotas allowed to allocate to a project in the routing group"`
	NamespaceScopeQuotaProportionCap float64 `json:"namespaceScopeQuotaProportionCap" pflag:",A floating point number between 0 and 1, specifying the maximum proportion of quotas allowed to allocate to a namespace in the routing group"`
}

type RateLimiter struct {
	Name         string          `json:"name" pflag:",The name of the rate limiter"`
	SyncPeriod   config.Duration `json:"syncPeriod" pflag:",The duration to wait before the cache is refreshed again"`
	Workers      int             `json:"workers" pflag:",Number of parallel workers to refresh the cache"`
	LruCacheSize int             `json:"lruCacheSize" pflag:",Size of the cache"`
	MetricScope  string          `json:"metricScope" pflag:",The prefix in Prometheus used to track metrics related to Presto"`
}

var (
	defaultConfig = Config{
		Environment:         URLMustParse(""),
		DefaultRoutingGroup: "adhoc",
		DefaultUser:         "flyte-default-user@lyft.com",
		RoutingGroupConfigs: []RoutingGroupConfig{{Name: "adhoc", Limit: 250}, {Name: "etl", Limit: 100}},
		RateLimiter: RateLimiter{
			Name:         "presto",
			SyncPeriod:   config.Duration{Duration: 3 * time.Second},
			Workers:      15,
			LruCacheSize: 2000,
			MetricScope:  "presto",
		},
	}

	prestoConfigSection = pluginsConfig.MustRegisterSubSection(prestoConfigSectionKey, &defaultConfig)
)

// Presto plugin configs
type Config struct {
	Environment         config.URL           `json:"environment" pflag:",Environment endpoint for Presto to use"`
	DefaultRoutingGroup string               `json:"defaultRoutingGroup" pflag:",Default Presto routing group"`
	DefaultUser         string               `json:"defaultUser" pflag:",Default Presto user"`
	RoutingGroupConfigs []RoutingGroupConfig `json:"routingGroupConfigs" pflag:"-,A list of cluster configs. Each of the configs corresponds to a service cluster"`
	RateLimiter         RateLimiter          `json:"rateLimiter" pflag:"Rate limiter config"`
}

// Retrieves the current config value or default.
func GetPrestoConfig() *Config {
	return prestoConfigSection.GetConfig().(*Config)
}

func SetPrestoConfig(cfg *Config) error {
	return prestoConfigSection.SetConfig(cfg)
}
