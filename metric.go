package razorpay

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Ref: https://godoc.org/github.com/prometheus/client_golang/prometheus/promhttp#ex-InstrumentRoundTripperDuration.

// metricCollector implements prometheus.Collector interface.
type metricCollector struct {
	inFlightGauge prometheus.Gauge
	counter       *prometheus.CounterVec
	dnsLatencyVec *prometheus.HistogramVec
	tlsLatencyVec *prometheus.HistogramVec
	histVec       *prometheus.HistogramVec
}

// NewPrometheusCollector configures and returns collector for http client.
func NewPrometheusCollector(client *http.Client, identifier string) prometheus.Collector {
	m := &metricCollector{}

	constLabels := map[string]string{"client_identifier": identifier}

	m.inFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "client_in_flight_requests",
		Help:        "A gauge of in-flight requests for the wrapped client.",
		ConstLabels: constLabels,
	})

	m.counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "client_api_requests_total",
			Help:        "A counter for requests from the wrapped client.",
			ConstLabels: constLabels,
		},
		[]string{"code", "method"},
	)

	// dnsLatencyVec uses custom buckets based on expected dns durations.
	// It has an instance label "event", which is set in the
	// DNSStart and DNSDonehook functions defined in the
	// InstrumentTrace struct below.
	m.dnsLatencyVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "dns_duration_seconds",
			Help:        "Trace dns latency histogram.",
			ConstLabels: constLabels,
			Buckets:     []float64{.005, .01, .025, .05},
		},
		[]string{"event"},
	)

	// tlsLatencyVec uses custom buckets based on expected tls durations.
	// It has an instance label "event", which is set in the
	// TLSHandshakeStart and TLSHandshakeDone hook functions defined in the
	// InstrumentTrace struct below.
	m.tlsLatencyVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "tls_duration_seconds",
			Help:        "Trace tls latency histogram.",
			ConstLabels: constLabels,
			Buckets:     []float64{.05, .1, .25, .5},
		},
		[]string{"event"},
	)

	// histVec has no labels, making it a zero-dimensional ObserverVec.
	m.histVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of request latencies.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)

	trace := &promhttp.InstrumentTrace{
		DNSStart: func(t float64) {
			m.dnsLatencyVec.WithLabelValues("dns_start").Observe(t)
		},
		DNSDone: func(t float64) {
			m.dnsLatencyVec.WithLabelValues("dns_done").Observe(t)
		},
		TLSHandshakeStart: func(t float64) {
			m.tlsLatencyVec.WithLabelValues("tls_handshake_start").Observe(t)
		},
		TLSHandshakeDone: func(t float64) {
			m.tlsLatencyVec.WithLabelValues("tls_handshake_done").Observe(t)
		},
	}

	// Instruments transport of client.
	client.Transport = promhttp.InstrumentRoundTripperInFlight(m.inFlightGauge,
		promhttp.InstrumentRoundTripperCounter(m.counter,
			promhttp.InstrumentRoundTripperTrace(trace,
				promhttp.InstrumentRoundTripperDuration(m.histVec, client.Transport),
			),
		),
	)

	return m
}

func (m *metricCollector) Describe(ch chan<- *prometheus.Desc) {
	m.inFlightGauge.Describe(ch)
	m.counter.Describe(ch)
	m.dnsLatencyVec.Describe(ch)
	m.tlsLatencyVec.Describe(ch)
	m.histVec.Describe(ch)
}

func (m *metricCollector) Collect(ch chan<- prometheus.Metric) {
	m.inFlightGauge.Collect(ch)
	m.counter.Collect(ch)
	m.dnsLatencyVec.Collect(ch)
	m.tlsLatencyVec.Collect(ch)
	m.histVec.Collect(ch)
}
