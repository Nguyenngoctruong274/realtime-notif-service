package config

import (
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

var (
	server           ServerCfg
	jaegerConfig     JaegerConfig
	dbCfg            DBCfg
	httpClient       HTTPClientCfg
	redisClient      RedisClientCfg
	prometheusConfig PrometheusCfg
	kafkaConfig      KafkaCfg
	keycloakConfig   KeycloakClientCfg
	postgresCfg      PostgresCfg
	websocketConfig  WebsocketCfg
)

type ServerCfg struct {
	ServerURL          string `envconfig:"SERVER_URL" default:"0.0.0.0"`
	LogLevel           string `envconfig:"LOG_LEVEL" default:"debug"`
	HTTPPort           int    `envconfig:"PORT" default:"8085"`
	Production         bool   `envconfig:"PRODUCTION" default:"false"`
	Env                string `envconfig:"ENVIRONMENT" default:"development"`
	LoggerSplunkURL    string `envconfig:"LOGGER_SPLUNK_URL" default:""`
	LoggerSplunkLayout string `envconfig:"LOGGER_SPLUNK_LAYOUT" default:""`
}

type DBCfg struct {
	UserName           string `envconfig:"USER_NAME" default:""`
	Password           string `envconfig:"PASSWORD" default:""`
	Host               string `envconfig:"HOST" default:""`
	DBName             string `envconfig:"DB_NAME" default:""`
	Port               string `envconfig:"PGPORT" default:""`
	SSLMode            string `envconfig:"SSL_MODE" default:""`
	SetMaxIdleConns    string `envconfig:"SET_MAX_IDLE" default:""`
	SetMaxOpenConns    string `envconfig:"SET_MAX_OPEN" default:""`
	SetConnMaxLifetime string `envconfig:"SET_CONN_MAX_LIFETIME" default:""`
}

type DataWarehouseCfg struct {
	UserName           string `envconfig:"WH_USER_NAME" default:""`
	Password           string `envconfig:"WH_PASSWORD" default:""`
	Host               string `envconfig:"WH_HOST" default:""`
	DBName             string `envconfig:"WH_DB_NAME" default:""`
	Port               string `envconfig:"WH_PGPORT" default:""`
	SSLMode            string `envconfig:"WH_SSL_MODE" default:""`
	SetMaxIdleConns    int    `envconfig:"WH_SET_MAX_IDLE" default:""`
	SetMaxOpenConns    int    `envconfig:"WH_SET_MAX_OPEN" default:""`
	SetConnMaxLifetime int    `envconfig:"WH_SET_CONN_MAX_LIFETIME" default:""`
}

type PostgresReportCfg struct {
	DBUrl              string `envconfig:"PG_URL_REPORT" default:""`
	SetMaxIdleConns    string `envconfig:"SET_MAX_IDLE" default:""`
	SetMaxOpenConns    string `envconfig:"SET_MAX_OPEN" default:""`
	SetConnMaxLifetime string `envconfig:"SET_CONN_MAX_LIFETIME" default:""`
}

type PostgresReplicateCfg struct {
	DBUrl              string `envconfig:"PG_URL_REPLICATE" default:""`
	SetMaxIdleConns    string `envconfig:"SET_MAX_IDLE" default:""`
	SetMaxOpenConns    string `envconfig:"SET_MAX_OPEN" default:""`
	SetConnMaxLifetime string `envconfig:"SET_CONN_MAX_LIFETIME" default:""`
}

type PostgresCDHCfg struct {
	DBUrl              string `envconfig:"PG_URL_CDH" default:""`
	SetMaxIdleConns    string `envconfig:"SET_MAX_IDLE" default:""`
	SetMaxOpenConns    string `envconfig:"SET_MAX_OPEN" default:""`
	SetConnMaxLifetime string `envconfig:"SET_CONN_MAX_LIFETIME" default:""`
}

type PostgresCfg struct {
	DBUrl              string `envconfig:"PG_URL" default:""`
	SetMaxIdleConns    string `envconfig:"SET_MAX_IDLE" default:""`
	SetMaxOpenConns    string `envconfig:"SET_MAX_OPEN" default:""`
	SetConnMaxLifetime string `envconfig:"SET_CONN_MAX_LIFETIME" default:""`
}

type HTTPClientCfg struct {
	ChatBotAPIURL string `envconfig:"CHAT_BOT_API_URL"`
	ChatBotAuth   string `envconfig:"CHAT_BOT_AUTH"`
	GGChatURL     string `envconfig:"GG_CHAT_URL"`
	ESAPIURL      string `envconfig:"ES_API_URL"`
}

type RedisClientCfg struct {
	RedisURL       string `envconfig:"REDIS_URL"`
	RedisSigleMode bool   `envconfig:"REDIS_SIGLE_MODE"`
}

type PrometheusCfg struct {
	PrometheusPort int    `envconfig:"PORT_METRIC"`
	ServiceName    string `envconfig:"PROM_SERVICE_NAME"`
}

type KafkaCfg struct {
	InitTopics        bool   `envconfig:"KAFKA_INIT_TOPIC"`
	Brokers           string `envconfig:"KAFKA_BROKERS"`
	GroupID           string `envconfig:"KAFKA_GROUP_ID"`
	PoolSize          int    `envconfig:"KAFKA_POOL_SIZE"`
	Partition         int    `envconfig:"KAFKA_PATITION"`
	ReplicationFactor int    `envconfig:"KAFKA_REPLICATION"`
	// kafkaTopics...
	TopicDLQ            string `envconfig:"KAFKA_TOPIC_DLQ"`
	TopicBudgetProfile  string `envconfig:"KAFKA_TOPIC_BUDGET_PROFILE"`
	TopicActivityLog    string `envconfig:"KAFKA_TOPIC_ACITIVITY"`
	TopicActivityEntity string `envconfig:"KAFKA_TOPIC_ACITIVITY_ENTITY"`
	// kafka retry opts...
	KafkaRetryAttempts     uint   `envconfig:"KAFKA_RETRY_ATTEMPTS"`
	KafkaRetryDelay        int    `envconfig:"KAFKA_RETRY_DELAYS"`
	PushFailedMessageToDLQ bool   `envconfig:"KAFKA_PUSH_FAILED_TO_DLQ"`
	DLQMessageKey          string `envconfig:"KAFKA_DLQ_MESSAGE_KEY"`

	// operation activity
	TopicOperationActivityLog string `envconfig:"KAFKA_TOPIC_OPERATION_ACTIVITY"`

	// notification
	TopicNotification    string `envconfig:"KAFKA_TOPIC_NOTIFICATION"`
	TopicNotificationDLQ string `envconfig:"KAFKA_TOPIC_NOTIFICATION_DLQ"`
}

type KeycloakClientCfg struct {
	KeycloakDomain    string `envconfig:"KEYCLOAK_DOMAIN"`
	KeycloakYSCDomain string `envconfig:"KEYCLOAK_YSC_DOMAIN"`
}

type WebsocketCfg struct {
	WebsocketSecretKey string `envconfig:"WEBSOCKET_SECRET_KEY"`
}

func (k *KafkaCfg) GetBrokers() []string {
	if len(kafkaConfig.Brokers) == 0 {
		return []string{}
	}
	return strings.Split(kafkaConfig.Brokers, ",")
}

func (k *KafkaCfg) GetTopicsConsume() []string {
	return []string{
		k.TopicBudgetProfile,
	}
}

func (k *KafkaCfg) GetTopicsNotification() []string {
	return []string{
		k.TopicNotification,
	}
}

type SenPayCfg struct {
	SenPayMerchantID             int    `envconfig:"SENPAY_MERCHANT_ID"`
	SenPayAccountMerchantID      int    `envconfig:"SENPAY_ACCOUNT_MERCHANT_ID"`
	SenPaySecurePassword         string `envconfig:"SENPAY_SECURE_PASS"`
	TopupFarmWalletDuration      int    `envconfig:"TOPUP_FARM_WALLET_DURATION"`
	LimitTopupFarmWalletDuration int    `envconfig:"LIMIT_TOPUP_FARM_WALLET_DURATION" default:"12"`
}

type PIMSCfg struct {
	PIMSDomain    string `envconfig:"PIMS_DOMAIN"`
	PIMSSecretKey string `envconfig:"PIMS_SECRET_KEY"`
}

func InitConfig() {
	configs := []interface{}{
		&server,
		&dbCfg,
		&jaegerConfig,
		&httpClient,
		&redisClient,
		&kafkaConfig,
		&keycloakConfig,
		&postgresCfg,
		&websocketConfig,
	}
	for _, instance := range configs {
		err := envconfig.Process("", instance)
		if err != nil {
			log.Fatalf("unable to init config: %v, err: %v", instance, err)
		}
	}
}

func ServerConfig() ServerCfg {
	return server
}

func DBConfig() DBCfg {
	return dbCfg
}

func HTTPClientConfig() HTTPClientCfg {
	return httpClient
}

func RedisConfig() RedisClientCfg {
	return redisClient
}

func PrometheusConfig() PrometheusCfg {
	return prometheusConfig
}

func KafkaConfig() KafkaCfg {
	return kafkaConfig
}

func KeycloakConfig() KeycloakClientCfg {
	return keycloakConfig
}

func PostgresConfig() PostgresCfg {
	return postgresCfg
}

func WebsocketConfig() WebsocketCfg {
	return websocketConfig
}
