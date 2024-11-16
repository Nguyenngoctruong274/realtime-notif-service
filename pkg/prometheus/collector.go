package prometheus

import (
	"yes4all/ads-noti-api/pkg/prometheus/prompusher"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector interface {
	Collect(collectors ...prometheus.Collector)
	UnCollect(collector prometheus.Collector)
}

type registryCollector struct {
	registerer prometheus.Registerer
}

func (r *registryCollector) Collect(collectors ...prometheus.Collector) {
	prometheus.MustRegister(collectors...)
}

func (r *registryCollector) UnCollect(collector prometheus.Collector) {
	prometheus.Unregister(collector)
}

func NewRegistryCollector() Collector {
	registerer := prometheus.NewRegistry()
	return &registryCollector{registerer: registerer}
}

type pusherCollector struct {
	pusher prompusher.Pusher
}

func NewPusherCollector(pusher prompusher.Pusher) Collector {
	return &pusherCollector{pusher: pusher}
}

func (r *pusherCollector) Collect(collectors ...prometheus.Collector) {
	for _, collector := range collectors {
		r.pusher.Collector(collector)
	}
}

func (r *pusherCollector) UnCollect(_ prometheus.Collector) {}
