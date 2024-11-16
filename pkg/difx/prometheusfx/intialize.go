package prometheusfx

import (
	"context"
	"fmt"
	"time"
	"yes4all/ads-noti-api/pkg/graceful"
	prometheus "yes4all/ads-noti-api/pkg/prometheus"
	"yes4all/ads-noti-api/pkg/prometheus/ginprometheus"
	"yes4all/ads-noti-api/pkg/prometheus/prompusher"
	"yes4all/ads-noti-api/pkg/prometheus/queueprometheus"
	"yes4all/ads-noti-api/pkg/xtime"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

const (
	stopTimeOut = time.Second
	waitTimeOut = time.Second
)

var Module = fx.Options(
	fx.Provide(providePromConfig, provideGinPrometheusMiddleware, provideQueuePrometheusMiddleware),
	fx.Invoke(initializeMetricsHandler),
)

func providePromConfig() prometheus.Config {
	return prometheus.ReadConfig()
}

func initializeMetricsHandler(lifecycle fx.Lifecycle, cfg prometheus.Config) {
	if !cfg.Enabled {
		fmt.Println("prometheus is disable")
		return
	}
	if !cfg.IsExposeMetricMode() {
		return
	}
	g := gin.New()
	g.Any("/metrics", MiddlewarePrometheus(), gin.WrapH(promhttp.Handler()))
	gracefulServer := graceful.NewService(graceful.WithStopTimeout(stopTimeOut), graceful.WithWaitTime(waitTimeOut))
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("expose metrics on port:", cfg.MetricsPort)
			go gracefulServer.StartServer(g, cfg.MetricsPort)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			gracefulServer.Close()
			return nil
		},
	})
}

func MiddlewarePrometheus() gin.HandlerFunc {
	metric := ginprometheus.GetPrometheusMetric()
	return func(c *gin.Context) {
		c.Next()
		fmt.Println("MiddlewarePrometheus")
		metric.ResetLatency()
		return
	}
}

func provideGinPrometheusMiddleware(lifecycle fx.Lifecycle, cfg prometheus.Config) ginprometheus.Middleware {
	if !cfg.Enabled {
		fmt.Println("prometheus is disable")
		return &ginprometheus.EmptyMiddleware{}
	}
	if cfg.IsExposeMetricMode() {
		return ginprometheus.NewMiddleware(cfg.ServiceName, prometheus.NewRegistryCollector())
	}
	pusher := initPrometheusPusher(lifecycle, cfg.GatewayConfig)
	return ginprometheus.NewMiddleware(cfg.ServiceName, prometheus.NewPusherCollector(pusher))
}

func provideQueuePrometheusMiddleware(lifecycle fx.Lifecycle, cfg prometheus.Config) queueprometheus.Middleware {
	if !cfg.Enabled {
		fmt.Println("prometheus is disable")
		return &queueprometheus.EmptyMiddleware{}
	}
	if cfg.IsExposeMetricMode() {
		return queueprometheus.NewMiddleware(cfg.ServiceName, prometheus.NewRegistryCollector())
	}
	pusher := initPrometheusPusher(lifecycle, cfg.GatewayConfig)
	return queueprometheus.NewMiddleware(cfg.ServiceName, prometheus.NewPusherCollector(pusher))
}

func initPrometheusPusher(lifecycle fx.Lifecycle, cfg prompusher.GatewayConfig) prompusher.Pusher {
	promPusher := prompusher.NewPrometheusPusher(cfg)
	cj := cron.New(cron.WithLocation(xtime.GMT07Location()))
	lifecycle.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			switch cfg.CronFunc {
			case prompusher.CronFnPush:
				err := promPusher.RegisterCronJobPush(cj, cfg.CronSpec)
				if err != nil {
					return err
				}
			case prompusher.CronFnAdd:
				err := promPusher.RegisterCronJobAdd(cj, cfg.CronSpec)
				if err != nil {
					return err
				}
			}
			cj.Start()
			return nil
		},
		OnStop: func(c context.Context) error {
			cj.Stop()
			return nil
		},
	})
	return promPusher
}
