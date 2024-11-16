package xprom

import (
	"context"
	"yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/utils/constants"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type Measurer interface {
	CounterInc(labels ...string)
	CounterIncError(err error, labels ...string)
	HistogramObserve(value float64, labels ...string)
}

var (
	counterLabels    = []string{eventLabel, actionLabel}
	counterMetricOpt = prometheus.CounterOpts{
		Namespace: constants.NameSpace,
		Subsystem: constants.PromSystem,
		Name:      "total_events",
		Help:      "Total events",
	}

	counterErrLabels    = []string{eventLabel, actionLabel}
	counterErrMetricOpt = prometheus.CounterOpts{
		Namespace: constants.NameSpace,
		Subsystem: constants.PromSystem,
		Name:      "total_error_events",
		Help:      "Total error events",
	}

	histogramLabels    = []string{eventLabel, actionLabel}
	histogramMetricOpt = prometheus.HistogramOpts{
		Namespace: constants.NameSpace,
		Subsystem: constants.PromSystem,
		Name:      "histogram",
		Help:      "histogram",
	}
)

type measurer struct {
	counterVec    *prometheus.CounterVec
	counterErrVec *prometheus.CounterVec
	histogramVec  *prometheus.HistogramVec
	opts          options
}

// CounterIncError increases error event with labels that's defined in counter vec.
// If labels is not match size, app's panic.
// Ex: define labels: []string{event, action}.
// We should call: CounterIncError("event_search_station","action_get_stations")
func (m *measurer) CounterIncError(err error, labels ...string) {
	if len(labels) != len(counterErrLabels) {
		logger.NewLogger().WithKeyword(context.Background(),
			"[PROMETHEUS] ignore inc counter error metrics coz labels are unmatched").
			WithFields(logrus.Fields{
				"expected": counterErrLabels,
				"actual":   labels,
			}).Warn()
		return
	}
	m.counterErrVec.WithLabelValues(labels...).Inc()
}

// CounterInc increases event with labels that's defined in counter vec. If labels is not match size, app's panic.
// Ex: define labels: []string{event, action}. We should call: CounterInc("event_search_station","action_get_stations")
func (m *measurer) CounterInc(labels ...string) {
	if len(labels) != len(counterLabels) {
		logger.NewLogger().WithKeyword(context.Background(),
			"[PROMETHEUS] ignore inc counter metrics coz labels are unmatched").
			WithFields(logrus.Fields{
				"expected": counterLabels,
				"actual":   labels,
			}).Warn()
	}
	m.counterVec.WithLabelValues(labels...).Inc()
}

// HistogramObserve increases event with labels that's defined in histogram vec.
// If labels is not match size, app's panic.
// Ex: define labels: []string{event, action}.
// We should call: HistogramObserve("event_search_station","action_get_stations")
func (m *measurer) HistogramObserve(value float64, labels ...string) {
	if len(labels) != len(histogramLabels) {
		logger.NewLogger().WithKeyword(context.Background(),
			"[PROMETHEUS] ignore observe histogram metrics coz labels are unmatched").
			WithFields(logrus.Fields{
				"expected": histogramLabels,
				"actual":   labels,
			}).Warn()
		return
	}
	m.histogramVec.WithLabelValues(labels...).Observe(value)
}

func NewMeasurer(opts ...Option) Measurer {
	optsArgs := NewDefPromOpt()
	for _, opt := range opts {
		opt.apply(&optsArgs)
	}
	m := initMeasurer(optsArgs)
	return m
}

func initMeasurer(opts options) *measurer {
	counterMetricOpt.ConstLabels = opts.ConstLabel
	counterVec := prometheus.NewCounterVec(counterMetricOpt, counterLabels)

	counterErrMetricOpt.ConstLabels = opts.ConstLabel
	counterErrVec := prometheus.NewCounterVec(counterErrMetricOpt, counterErrLabels)

	histogramMetricOpt.ConstLabels = opts.ConstLabel
	histogramVec := prometheus.NewHistogramVec(histogramMetricOpt, histogramLabels)

	opts.Register.MustRegister(counterVec, counterErrVec, histogramVec)
	return &measurer{opts: opts, counterVec: counterVec, counterErrVec: counterErrVec, histogramVec: histogramVec}
}
