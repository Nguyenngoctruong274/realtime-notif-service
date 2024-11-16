package metricsfx

import (
	"go.uber.org/fx"

	"yes4all/ads-noti-api/pkg/xprom"
)

var Module = fx.Provide(
	provideMeasurerProvider)

func provideMeasurerProvider() xprom.MeasurerProvider {
	provider := xprom.NewDefaultProvider()
	return provider
}
