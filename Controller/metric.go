// controller/metrics.go

package controller

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	addOrder    prometheus.Counter
	cancelOrder prometheus.Counter
	matchOrder  prometheus.Counter
}

func newMetrics(gatherer ametrics.MultiGatherer) (*Metrics, error) {
	m := &Metrics{
		addOrder: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orderbook_add_order_total",
			Help: "Total number of AddOrder actions executed",
		}),
		cancelOrder: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orderbook_cancel_order_total",
			Help: "Total number of CancelOrder actions executed",
		}),
		matchOrder: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orderbook_match_order_total",
			Help: "Total number of MatchOrder actions executed",
		}),
	}

	// Register metrics
	registry := prometheus.NewRegistry()
	err := registry.Register(m.addOrder)
	if err != nil {
		return nil, err
	}
	err = registry.Register(m.cancelOrder)
	if err != nil {
		return nil, err
	}
	err = registry.Register(m.matchOrder)
	if err != nil {
		return nil, err
	}

	// Add registry to the gatherer
	gatherer.Register("orderbook", registry)

	return m, nil
}
