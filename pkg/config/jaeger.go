package config

type JaegerConfig struct {
	Enabled      bool    `envconfig:"JAEGER_ENABLED" default:"false"`
	ServiceName  string  `envconfig:"JAEGER_SERVICE_NAME" default:"red-package-api"`
	Endpoint     string  `envconfig:"JAEGER_ENDPOINT" default:"127.0.0.1:6382"`
	SamplerType  string  `envconfig:"JAEGER_SAMPLER_TYPE" default:"const"`
	SamplerParam float64 `envconfig:"JAEGER_SAMPLER_PARAM" default:"1"`
}

func TracingConfig() JaegerConfig {
	return jaegerConfig
}
