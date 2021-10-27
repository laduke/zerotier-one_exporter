package domain_test

import (
	"testing"
	"github.com/laduke/zerotier-one_exporter/internal/domain"
)

func TestNewMetrics(t *testing.T) {
	var latency  map[string]float64
	metrics := domain.NewMetrics(0, 0, latency)

	if metrics.TotalPeers != 0 {
		t.Errorf("got %f want %f", metrics.TotalPeers, 0.0)
	}

	if metrics.DirectPeers != 0 {
		t.Errorf("got %f want %f", metrics.DirectPeers, 0.0)
	}

}
