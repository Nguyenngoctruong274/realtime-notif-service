package prometheus

import (
	"yes4all/ads-noti-api/pkg/env"
	"yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/prometheus/prompusher"

	"github.com/spf13/viper"
)

const (
	modeExposeMetrics  = "expose_metrics"
	modePushGateway    = "push_gateway"
	defaultServiceName = "yams_service"
	defaultMetricsPort = "9090"
)

var DefaultPromConfig = Config{
	Enabled:       false,
	ServiceName:   defaultServiceName,
	Mode:          modeExposeMetrics,
	MetricsPort:   defaultMetricsPort,
	GatewayConfig: prompusher.DefaultPromGatewayConfig,
}

type Config struct {
	Enabled       bool
	ServiceName   string
	Mode          string
	MetricsPort   string
	GatewayConfig prompusher.GatewayConfig
}

func (c Config) IsExposeMetricMode() bool {
	return c.Mode == modeExposeMetrics
}

func ReadConfig() Config {
	cfg := DefaultPromConfig
	cfg.Enabled = viper.GetBool(env.PromEnable)
	serviceName := viper.GetString(env.PromServiceName)
	if serviceName != "" {
		cfg.ServiceName = serviceName
	}
	mode := viper.GetString(env.PromMode)
	if isValidMode(mode) {
		cfg.Mode = mode
	} else {
		logger.NewLogger().Warn("prometheus mode is invalid or empty, it should be `expose_metrics` or `push_gateway`," +
			" use default mode `expose_metrics`")
	}
	metricsPort := viper.GetString(env.PromMetricsPort)
	if metricsPort != "" {
		cfg.MetricsPort = metricsPort
	}
	cfg.GatewayConfig = prompusher.ReadGatewayConfig()
	logger.NewLogger().Debugf("prometheus config: %+v", cfg)
	return cfg
}

func isValidMode(mode string) bool {
	return mode == modeExposeMetrics || mode == modePushGateway
}
