package ztone

import (
	one "github.com/erikh/go-ztone"
	"github.com/laduke/zerotier-one_exporter/internal/domain"
)

type ztone struct {
	*one.Client
}

func New(token string) *ztone {
	return &ztone{one.NewClient(token)}
}

func (repo *ztone) Get() (domain.Metrics, error) {
	res, err := repo.Peers()
	if err != nil {
		return domain.Metrics{}, err
	}

	metrics := domain.Metrics{}
	metrics.TotalPeers = float64(len(res))
	metrics.Latency = make(map[string]float64, len(res))

	metrics.DirectPeers = 0
	for _, v := range res {
		if IsDirect(*v) {
			metrics.DirectPeers = metrics.DirectPeers + 1
		}

		metrics.Latency[v.Address] = Latency(*v)
	}

	return metrics, nil
}

func IsDirect(p one.Peer) bool {
	direct := false
	for _, v := range p.Paths {
		if v.Preferred {
			direct = true
			break
		}
	}
	return direct
}

func Latency(p one.Peer) float64 {
	var l float64
	if p.Latency != -1 {
		l = float64(p.Latency) / 1000
	} else {
		l = -1
	}
	return l
}
