package env

const (
	Production = "PRODUCTION" // true - fasle
	Port       = "PORT"
	// MONGODB ENV
	MongoURI    = "MONGO_URI"
	MongoDBName = "MONGO_DB"
	// PROMETHEUS ENV
	PromEnable                = "PROM_ENABLE"
	PromServiceName           = "PROM_SERVICE_NAME"
	PromMode                  = "PROM_MODE"
	PromMetricsPort           = "PROM_METRICS_PORT"
	PromURL                   = "PROM_URL"
	PromJobName               = "PROM_JOBNAME"
	PromInstanceID            = "PROM_INSTANCE_ID"
	PromRequestInterval       = "PROM_REQUEST_INTERVAL"
	PromHTTPTimeOut           = "PROM_HTTP_TIMEOUT"
	PromHystrixTimeout        = "PROM_HYSTRIX_TIMEOUT"
	PromSleepWindow           = "PROM_SLEEP_WINDOW" // #nosec
	PromMaxConcurrentRequests = "PROM_MAX_CONCURRENT_REQUESTS"
	PromErrorPercentThreshold = "PROM_ERROR_PERCENT_THRESHOLD"
	PromCronFunc              = "PROM_CRON_FUNC"
	PromCronSpec              = "PROM_CRON_SPEC"

	// APM ENV
	APMServerURL             = "ELASTIC_APM_SERVER_URL"
	APMServerURLs            = "ELASTIC_APM_SERVER_URLS"
	APMServiceName           = "ELASTIC_APM_SERVICE_NAME"
	APMServiceVersion        = "ELASTIC_APM_SERVICE_VERSION"
	APMEnvironment           = "ELASTIC_APM_ENVIRONMENT"
	APMTransactionSampleRate = "ELASTIC_APM_TRANSACTION_SAMPLE_RATE"
	// REDIS ENV
	RedisURL        = "REDIS_URL"
	RedisPwd        = "REDIS_PWD"
	RedisSingleMode = "REDIS_SINGLE_MODE"
)
