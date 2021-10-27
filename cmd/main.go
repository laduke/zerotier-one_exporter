package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/laduke/zerotier-one_exporter/internal/authtokenservice"
	"github.com/laduke/zerotier-one_exporter/internal/promexporter"
	"github.com/laduke/zerotier-one_exporter/internal/ztoneservice"
)

// TODO make a config service, configure http port # and authtoken and flag to collect addresses or not

func main() {
	token := authtoken.Guess()
	oneService := ztone.New(token)

	reg := prometheus.NewPedanticRegistry()
	exporter := promexporter.New(oneService)
	reg.MustRegister(exporter)

	fmt.Println("starting")

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":1971", nil)
}
