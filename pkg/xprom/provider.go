package xprom

import "errors"

type MeasurerProvider interface {
	Provide(biz Biz) Measurer
	Add(biz Biz, measurer Measurer) error
}

type measurerProvider struct {
	measurerMap map[Biz]Measurer
}

func (m *measurerProvider) Add(biz Biz, measurer Measurer) error {
	_, ok := m.measurerMap[biz]
	if ok {
		return errors.New("measurer is existed")
	}
	m.measurerMap[biz] = measurer
	return nil
}

func (m *measurerProvider) Provide(biz Biz) Measurer {
	measurer, ok := m.measurerMap[biz]
	if ok {
		return measurer
	}
	return m.measurerMap[BizUnknown]
}

func NewDefaultProvider() MeasurerProvider {
	return &measurerProvider{
		measurerMap: map[Biz]Measurer{
			BizUnknown: NewMeasurer(WithBiz(BizUnknown)),
			BizGeneral: NewMeasurer(WithBiz(BizGeneral)),
		},
	}
}

func NewProvider(m map[Biz]Measurer) MeasurerProvider {
	if m == nil {
		m = map[Biz]Measurer{
			BizUnknown: NewMeasurer(WithBiz(BizUnknown)),
			BizGeneral: NewMeasurer(WithBiz(BizGeneral)),
		}
	}
	return &measurerProvider{
		measurerMap: m,
	}
}
