package prompusher

import (
	"fmt"
	"time"
	"yes4all/ads-noti-api/pkg/env"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

const (
	CronFnPush = "push"
	CronFnAdd  = "add"
)
const (
	defaultHTTPTimeOut           = 500 * time.Millisecond
	defaultHystrixTimeout        = 200 * time.Millisecond
	defaultSleepWindow           = 1000 * time.Millisecond
	defaultMaxConcurrentRequests = 50
	defaultErrorPercentThreshold = 10
	defaultRequestInterval       = 300 * time.Millisecond
	defaultCronSpec              = "@every 5s"
	defaultGatewayURL            = "http://localhost:9091"
	defaultJobName               = "job-default"
	defaultInstanceID            = "instance-id-default"
)

var defaultHystrixConfig = HystrixConfig{
	HTTPTimeout:           defaultHTTPTimeOut,
	HystrixTimeout:        defaultHystrixTimeout,
	SleepWindow:           defaultSleepWindow,
	MaxConcurrentRequests: defaultMaxConcurrentRequests,
	ErrorPercentThreshold: defaultErrorPercentThreshold,
	RequestInterval:       defaultRequestInterval,
}

var defaultCronJobConfig = CronJobConfig{
	CronFunc: CronFnPush,
	CronSpec: defaultCronSpec,
}
var DefaultPromGatewayConfig = GatewayConfig{
	GatewayURL:    defaultGatewayURL,
	InstanceID:    defaultInstanceID,
	JobName:       defaultJobName,
	HystrixConfig: defaultHystrixConfig,
	CronJobConfig: defaultCronJobConfig,
}

type GatewayConfig struct {
	GatewayURL string
	InstanceID string
	JobName    string
	HystrixConfig
	CronJobConfig
}

type HystrixConfig struct {
	HTTPTimeout           time.Duration
	HystrixTimeout        time.Duration
	SleepWindow           time.Duration
	MaxConcurrentRequests int
	ErrorPercentThreshold int
	RequestInterval       time.Duration
}
type CronJobConfig struct {
	CronFunc string
	CronSpec string
}

func ReadGatewayConfig() GatewayConfig {
	cfg := DefaultPromGatewayConfig
	gatewayURL := viper.GetString(env.PromURL)
	if gatewayURL != "" {
		cfg.GatewayURL = gatewayURL
	}
	jobName := viper.GetString(env.PromJobName)
	if jobName != "" {
		cfg.JobName = jobName
	}
	requestInterval := viper.GetDuration(env.PromRequestInterval)
	if requestInterval != 0 {
		cfg.RequestInterval = requestInterval
	}
	httpTimeout := viper.GetDuration(env.PromHTTPTimeOut)
	if httpTimeout != 0 {
		cfg.HTTPTimeout = httpTimeout
	}
	hystrixTimeout := viper.GetDuration(env.PromHystrixTimeout)
	if hystrixTimeout != 0 {
		cfg.HystrixTimeout = hystrixTimeout
	}
	sleepWindow := viper.GetDuration(env.PromSleepWindow)
	if sleepWindow != 0 {
		cfg.SleepWindow = sleepWindow
	}
	maxConcurrentRequests := viper.GetInt(env.PromMaxConcurrentRequests)
	if maxConcurrentRequests != 0 {
		cfg.MaxConcurrentRequests = maxConcurrentRequests
	}
	errorPercentThreshold := viper.GetInt(env.PromErrorPercentThreshold)
	if errorPercentThreshold != 0 {
		cfg.ErrorPercentThreshold = errorPercentThreshold
	}
	cronFunc := viper.GetString(env.PromCronFunc)
	if isValidCronFunc(cronFunc) {
		cfg.CronFunc = cronFunc
	}
	cronSpec := viper.GetString(env.PromCronSpec)
	if _, err := cron.ParseStandard(cronSpec); err == nil {
		cfg.CronSpec = cronSpec
	}

	instanceID := viper.GetString(env.PromInstanceID)
	if len(instanceID) == 0 {
		instanceID = fmt.Sprint(time.Now().Format("060102_15:04"))
	}
	cfg.InstanceID = instanceID
	return cfg
}

func isValidCronFunc(fn string) bool {
	return fn == CronFnPush || fn == CronFnAdd
}
