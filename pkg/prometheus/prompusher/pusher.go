package prompusher

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/robfig/cron/v3"

	"yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/xerrors"
)

//go:generate mockgen -source=pusher.go -package=pushermock -destination=gomock/pusher.go

type Pusher interface {
	Collector(c prometheus.Collector)
	AddCounterVec(opts prometheus.CounterOpts, labels []string) *prometheus.CounterVec
	AddHistogramVec(opts prometheus.HistogramOpts, labels []string) *prometheus.HistogramVec
	AddGaugeVec(opts prometheus.GaugeOpts, labels []string) *prometheus.GaugeVec
	AddSummaryVec(opts prometheus.SummaryOpts, labels []string) *prometheus.SummaryVec
}

type prometheusPusher struct {
	pusher *push.Pusher
}

func NewPrometheusPusher(conf GatewayConfig) *prometheusPusher {
	pusher := push.New(conf.GatewayURL,
		fmt.Sprintf("%s_%s", conf.JobName, conf.InstanceID))

	client := NewHTTPClient(conf)
	pusher.Client(client)

	pusher.Collector(promLatencyCollector)
	pusher.Collector(promRequestBodySizeCollector)
	pusher.Collector(promRequestCounterCollector)

	promPusher := &prometheusPusher{pusher}
	return promPusher
}

func (p *prometheusPusher) AddCounterVec(opts prometheus.CounterOpts, labels []string) *prometheus.CounterVec {
	collector := prometheus.NewCounterVec(opts, labels)
	p.pusher.Collector(collector)
	return collector
}

func (p *prometheusPusher) AddHistogramVec(opts prometheus.HistogramOpts, labels []string) *prometheus.HistogramVec {
	collector := prometheus.NewHistogramVec(opts, labels)
	p.pusher.Collector(collector)
	return collector
}

func (p *prometheusPusher) AddGaugeVec(opts prometheus.GaugeOpts, labels []string) *prometheus.GaugeVec {
	collector := prometheus.NewGaugeVec(opts, labels)
	p.pusher.Collector(collector)
	return collector
}

func (p *prometheusPusher) AddSummaryVec(opts prometheus.SummaryOpts, labels []string) *prometheus.SummaryVec {
	collector := prometheus.NewSummaryVec(opts, labels)
	p.pusher.Collector(collector)
	return collector
}

func (p *prometheusPusher) Collector(c prometheus.Collector) {
	p.pusher.Collector(c)
}

func (p *prometheusPusher) push() {
	err := p.pusher.Push()

	if err != nil && !xerrors.Is(err, xerrors.Success) {
		logger.NewLogger().Error("failed to push prom, error: ", err)
	}
}

func (p *prometheusPusher) add() {
	err := p.pusher.Add()

	if err != nil && xerrors.Is(err, xerrors.Success) {
		logger.NewLogger().Error("failed to add prom, error: ", err)
	}
}

func (p *prometheusPusher) RegisterCronJobPush(cron *cron.Cron, spec string) error {
	_, err := cron.AddFunc(spec, p.push)
	return err
}

func (p *prometheusPusher) RegisterCronJobAdd(cron *cron.Cron, spec string) error {
	_, err := cron.AddFunc(spec, p.add)
	return err
}
