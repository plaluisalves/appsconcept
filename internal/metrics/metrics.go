package metrics

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Service defines the interface for the metrics service.
type Service interface {
	Data() ([]MetricRow, []MetricRow, error)
	IncHits(values ...any)
}

type Metrics struct {

	// control metrics
	usePrometheus bool
	useCustom     bool

	// custom metrics
	metricsMap map[string]float64
	mu         sync.RWMutex

	// prometheus metrics
	FizzBuzzCounterRequests *prometheus.CounterVec // will collect metrics for FizzBuzz requests in memory
}

type MetricRow struct {
	Request string
	Total   float64
}

// Fixed order of label names
var fixedOrder = []string{"int1", "int2", "limit", "str1", "str2"}

// New creates a new metrics service.
func New(opts ...Option) Service {

	m := &Metrics{}
	for _, opt := range opts {
		opt(m)
	}

	switch {
	case m.useCustom && m.usePrometheus:
		panic("cannot use both custom and prometheus metrics at the same time!")
	case m.useCustom:
		m.metricsMap = make(map[string]float64)
	case m.usePrometheus:
		m.FizzBuzzCounterRequests = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fizzbuzz_requests_total",
				Help: "Total number of FizzBuzz requests",
			},
			fixedOrder, // []string{"int1", "int2", "limit", "str1", "str2"},
		)

		prometheus.MustRegister(m.FizzBuzzCounterRequests)

	default:
		// No metrics
		log.Println("metrics are disabled.")
	}

	return m

}

type Option func(*Metrics)

func WithPrometheus() Option {
	return func(m *Metrics) {
		m.usePrometheus = true
	}
}

func WithCustomSolution() Option {
	return func(m *Metrics) {
		m.useCustom = true
	}
}

// Data returns all collected metrics and the ones with the highest total.
func (m *Metrics) Data() ([]MetricRow, []MetricRow, error) {
	var rows []MetricRow
	var maxs []MetricRow
	var maxTotal float64

	//  track max
	addRow := func(r MetricRow) {
		rows = append(rows, r)
		if r.Total > maxTotal {
			maxTotal = r.Total
			maxs = []MetricRow{r} // new max found, reset list
		} else if r.Total == maxTotal {
			maxs = append(maxs, r) // same max, append
		}
	}

	if m.usePrometheus {
		metrics, err := prometheus.DefaultGatherer.Gather()
		if err != nil {
			return rows, maxs, fmt.Errorf("failed to gather prometheus metrics: %w", err)
		}

		for _, mf := range metrics {
			if mf.GetName() == "fizzbuzz_requests_total" {
				for _, metric := range mf.GetMetric() {
					labels := []string{}
					for _, lp := range metric.GetLabel() {
						labels = append(labels, lp.GetName()+"="+lp.GetValue())
					}
					addRow(MetricRow{
						Request: strings.Join(labels, "&"),
						Total:   metric.GetCounter().GetValue(),
					})
				}
			}
		}
	}

	if m.useCustom {
		m.mu.Lock()
		defer m.mu.Unlock()
		for k, v := range m.metricsMap {
			addRow(MetricRow{
				Request: k,
				Total:   v,
			})
		}
	}

	return rows, maxs, nil
}

// IncHits increments the hit counter for a specific combination of query parameters.
func (m *Metrics) IncHits(values ...any) {

	// Convert all values to string
	lvs := make([]string, len(values))
	for i, v := range values {
		lvs[i] = fmt.Sprint(v)
	}

	if m.usePrometheus {
		m.FizzBuzzCounterRequests.WithLabelValues(lvs...).Inc()
	}

	if m.useCustom {
		m.mu.Lock()
		defer m.mu.Unlock()

		// Build the key in the format field=value&...
		pairs := make([]string, len(lvs))
		for i, v := range lvs {
			// Prevent panic if fewer values than fieldNames are passed
			name := fmt.Sprintf("field%d", i)
			if i < len(fixedOrder) {
				name = fixedOrder[i]
			}
			pairs[i] = name + "=" + v
		}

		key := strings.Join(pairs, "&")
		m.metricsMap[key]++
	}
}
