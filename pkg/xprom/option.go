package xprom

import (
	"github.com/prometheus/client_golang/prometheus"
)

type options struct {
	Biz        Biz
	ConstLabel map[string]string
	Register   prometheus.Registerer
}

func NewDefPromOpt() options {
	return options{
		ConstLabel: map[string]string{},
		Register:   prometheus.DefaultRegisterer,
	}
}

func NewBasePromOpt(biz Biz) options {
	return options{
		Biz: biz,
		ConstLabel: map[string]string{
			bizLabel: string(biz),
		},
		Register: prometheus.DefaultRegisterer,
	}
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(args *options) {
	f(args)
}

func WithBiz(biz Biz) Option {
	return optionFunc(func(args *options) {
		args.ConstLabel[bizLabel] = string(biz)
	})
}

func WithRegister(register prometheus.Registerer) Option {
	return optionFunc(func(args *options) {
		args.Register = register
	})
}

func WithConstLabels(labels map[string]string) Option {
	return optionFunc(func(args *options) {
		if args.ConstLabel == nil {
			args.ConstLabel = make(map[string]string)
		}
		for label, value := range labels {
			args.ConstLabel[label] = value
		}
	})
}
