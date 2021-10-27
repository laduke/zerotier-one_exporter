package promexporter

import (
	"github.com/laduke/zerotier-one_exporter/internal/domain"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "zerotier"
)

type Exporter struct {
	metricsService domain.ZeroTierService
}

func New(srv domain.ZeroTierService) *Exporter {
	return &Exporter{
		metricsService: srv,
	}
}

var totalPeers = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: namespace,
	Name:      "peers_total",
	Help:      "count of peers",
})

var peers = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: namespace,
	Name:      "peers",
	Help:      "count of peers",
}, []string{"connection"})

var latency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: namespace,
	Name:      "latency",
	Help:      "latency",
}, []string{"address"})

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- peers.WithLabelValues("connection").Desc()
	ch <- totalPeers.Desc()
	ch <- latency.WithLabelValues("latency").Desc()
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	metrics, err := e.metricsService.Get()

	if err != nil {
		// TODO
		// log error
	} else {
		ch <- prometheus.MustNewConstMetric(totalPeers.Desc(), prometheus.GaugeValue, metrics.TotalPeers)
		ch <- prometheus.MustNewConstMetric(peers.WithLabelValues("connection").Desc(), prometheus.GaugeValue, metrics.DirectPeers, "direct")

		for k, v := range metrics.Latency {
			ch <- prometheus.MustNewConstMetric(latency.WithLabelValues("address").Desc(), prometheus.GaugeValue, v, k)
		}
	}

}
