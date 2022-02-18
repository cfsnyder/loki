package journal

import "github.com/prometheus/client_golang/prometheus"

// Metrics holds a set of journald metrics.
type Metrics struct {
	reg prometheus.Registerer

	journalFollowFailures prometheus.Counter
}

// NewMetrics creates a new set of journald metrics. If reg is non-nil, the
// metrics will be registered.
func NewMetrics(reg prometheus.Registerer) *Metrics {
	var m Metrics
	m.reg = reg

	m.journalFollowFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "promtail",
		Name:      "journal_follow_failure_total",
		Help:      "Total number of times that the journal position had to be reset because journal iteration failed. This usually happens when promtail is following too slowly and journald deletes entries before they're processed. Likely indicates that log lines were missed.",
	})

	if reg != nil {
		reg.MustRegister(
			m.journalFollowFailures,
		)
	}

	return &m
}
