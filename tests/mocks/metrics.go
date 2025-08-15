package mocks

import "appsceoncept/internal/metrics"

/*******MOCK METRICS SERVICE *******/
type MockMetricsService struct {
	DataFunc    func() ([]metrics.MetricRow, []metrics.MetricRow, error)
	IncHitsFunc func(values ...any)
}

func (m *MockMetricsService) Data() ([]metrics.MetricRow, []metrics.MetricRow, error) {
	if m.DataFunc != nil {
		return m.DataFunc()
	}
	return nil, nil, nil
}

func (m *MockMetricsService) IncHits(values ...any) {
	if m.IncHitsFunc != nil {
		m.IncHitsFunc(values...)
	}
}
