package domain

import ()

type Metrics struct {
	TotalPeers  float64
	DirectPeers float64
	Latency     map[string]float64
}

func NewMetrics(peerCount float64, directCount float64, latency map[string]float64) Metrics {
	return Metrics{TotalPeers: peerCount, DirectPeers: directCount, Latency: latency}
}

// interface

type ZeroTierService interface {
	Get() (Metrics, error)
}
