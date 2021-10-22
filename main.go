package main

import (
	// "fmt"
	"io/ioutil"
	"net/http"
	"os"

	one "github.com/erikh/go-ztone"
	"github.com/laduke/zerotier-one_exporter/authtoken"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	reg := prometheus.NewPedanticRegistry()

	token := os.Getenv("ZEROTIER_ONE_TOKEN")
	if token == "" {
		b, err := ioutil.ReadFile(authtoken.TokenPath())
		if err != nil { panic(err) }

		token = string(b)
	}

	c := one.NewClient(token)

	d := MyClient{ client: *c }


	exporter := NewExporter(&d)

	reg.MustRegister(exporter)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":1971", nil)

	// fmt.Printf("%v+", res)
}

type Noder interface {
	GetPeers() ([]MetricPeer, error)
	GetStatus() (MetricStatus, error)
	GetNetworks() ([]MetricNetwork, error)
}

type Exporter struct {
	client Noder
}

type MyClient struct {
	client one.Client
}

func NewExporter(client Noder) *Exporter {
	return &Exporter{
		client: client,
	}
}

func (c *MyClient) GetPeers() ([]MetricPeer, error) {
	res, err := c.client.Peers()
	if err != nil {
		return []MetricPeer{}, err
	}

	peers := make([]MetricPeer, len(res))

	for i, p := range res {
		peer := PeerToMetricPeer(p)
		for _, path := range p.Paths{
			path2 := PeerPathToMetricPeerPath(&path)
			peer.Paths = append(peer.Paths, path2)
		}
		peers[i] = peer
	}

	return peers, nil
}

func (c *MyClient) GetStatus() (MetricStatus, error) {
	res, err := c.client.Status()
	if err != nil {
		return MetricStatus{}, err
	}

	status := StatusToMetricStatus(res)
	return status, nil
}

func (c *MyClient) GetNetworks() ([]MetricNetwork, error) {
	res, err := c.client.Networks()
	if err != nil {
		return []MetricNetwork{}, err
	}

	networks := make([]MetricNetwork, len(res))

	for i, p := range res {
		network := NetworkToMetricNetwork(p)
		networks[i] = network
	}

	return networks, nil
}

const (
	namespace = "zerotier"
)

var peerRoles = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: namespace,
	Name: "peer_roles",
	Help: "count of peers by role",
}, []string{"role"})

var peerConn = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: namespace,
	Name: "peer_connections",
	Help: "count of peers by connection",
}, []string{"connection"})

var peerLatency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: namespace,
	Name: "peer_latency",
	Help: "count of peers by connection",
}, []string{"address", "role", "version"})

var statusGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: namespace,
	Name: "node_status",
	Help: "online or offline. (Can talk to roots)",
}, []string{"version", "address"})

var networkStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: namespace,
	Name: "network",
	Help: "network",
}, []string{"status"})

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	peers, _ := e.client.GetPeers()

	direct := CountPeerConnections(&peers)

	status, _ := e.client.GetStatus()

	networks, _ := e.client.GetNetworks()
	networkCounts := CountNetworks(&networks)

	planets, leafs, moons := CountPeerRoles(&peers)

	ch <- prometheus.MustNewConstMetric(statusGauge.WithLabelValues("version", "address").Desc(), prometheus.GaugeValue, status.OnlineFloat(), status.Version, status.Address)

	ch <- prometheus.MustNewConstMetric(peerRoles.WithLabelValues("role").Desc(), prometheus.GaugeValue, leafs, "leaf")
	ch <- prometheus.MustNewConstMetric(peerRoles.WithLabelValues("role").Desc(), prometheus.GaugeValue, planets, "planet")
	ch <- prometheus.MustNewConstMetric(peerRoles.WithLabelValues("role").Desc(), prometheus.GaugeValue, moons, "moon")

	ch <- prometheus.MustNewConstMetric(peerConn.WithLabelValues("connection").Desc(), prometheus.GaugeValue, direct, "direct")
	ch <- prometheus.MustNewConstMetric(peerConn.WithLabelValues("connection").Desc(), prometheus.GaugeValue, float64(len(peers)) - direct, "relay")

	for k, v := range(networkCounts) {
		ch <- prometheus.MustNewConstMetric(networkStatus.WithLabelValues("status").Desc(),prometheus.GaugeValue, v, k)
	}

	// TODO make enabling configurable maybe
	for _, v := range peers {
		ch <- prometheus.MustNewConstMetric(peerLatency.WithLabelValues("address", "role", "version").Desc(), prometheus.GaugeValue, float64(v.Latency), v.Address, v.Role, v.Version)
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- peerRoles.WithLabelValues("role").Desc()
	ch <- peerConn.WithLabelValues("connection").Desc()
	ch <- peerLatency.WithLabelValues("address", "role", "version").Desc()

	ch <- statusGauge.WithLabelValues("version", "address").Desc()

	ch <- networkStatus.WithLabelValues("status").Desc()
}

type MetricPeer struct {
	Address string
	Latency float64
	Paths []MetricPeerPath
	Role string
	Version string
}

type MetricPeerPath struct {
	Active bool
	Address string
	Expired bool
	Preferred bool
}

type MetricStatus struct {
	Address string
	Online bool
	Version string
	TCPFallbackActive bool
}

type MetricNetwork struct {
	ID string
	Status string
}

func StatusToMetricStatus (old *one.NetworkStatus) MetricStatus {
	var status MetricStatus
	status.Address = old.Address
	status.Online = old.Online
	status.Version = old.Version
	status.TCPFallbackActive = old.TCPFallbackActive

	return status
}

func PeerToMetricPeer (peer *one.Peer) MetricPeer {
	var peer2 MetricPeer
	peer2.Latency = float64(peer.Latency)
	peer2.Address = peer.Address
	peer2.Role = peer.Role
	peer2.Version = peer.Version

	return peer2
}

func NetworkToMetricNetwork (old *one.Network) MetricNetwork {
	var network2 MetricNetwork
	network2.ID = old.ID
	network2.Status = old.Status
	return network2
}

func PeerPathToMetricPeerPath (path *one.Path) MetricPeerPath {
	var path2 MetricPeerPath
	path2.Active = path.Active
	path2.Preferred = path.Preferred

	return path2
}

func CountPeerRoles(peers *[]MetricPeer) (float64, float64, float64) {
	leafs := 0
	planets := 0
	moons := 0

	for _, v := range *peers {
		if v.IsLeaf()  {
			leafs = leafs + 1
		} else if v.IsRoot() {
			planets = planets + 1
		} else if v.IsMoon() {
			moons = moons + 1
		}
	}

	return float64(planets), float64(leafs), float64(moons)
}

func CountPeerConnections(peers *[]MetricPeer) (float64) {
	direct := 0

	for _, v := range *peers {
		if v.IsDirect() {
			direct++
		}
	}

	return float64(direct)
}

func CountNetworks(networks *[]MetricNetwork) (map[string]float64) {
	counts := make(map[string]float64)

	for _, v := range *networks {
		status := v.Status
		curr := counts[status]
		counts[status] = curr + 1
	}

	return counts
}

func (p MetricPeer) IsDirect() bool {
	direct := false
	for _, v := range p.Paths {
		if v.Preferred == true {
			direct =  true
			break
		}
	}
	return direct
	// return len(p.Paths) > 0
}

func (p MetricPeer) IsController(networks *[]MetricNetwork) bool {
	res := false
	for _, v := range *networks {
		if (v.ID[0:10] == p.Address) {
			res = true
			break
		}
	}
	return res
}

func (p MetricPeer) IsLeaf() bool {
	return p.Role == "LEAF"
}

func (p MetricPeer) IsRoot() bool {
	return p.Role == "PLANET"
}

func (p MetricPeer) IsMoon() bool {
	return p.Role == "MOON"
}

func (p MetricPeer) Conn() string {
	if p.IsDirect() {
		return "DIRECT"
	} else { return "RELAY" }
}

func (n MetricStatus) OnlineFloat() float64 {
	if n.Online {
		return 1
	} else { return 0 }
}
