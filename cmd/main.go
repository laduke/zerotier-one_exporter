package main

import (
	"io/ioutil"
	"net/http"
	"os"

	// "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/laduke/zerotier-one_exporter/internal/authtokenservice"
	"github.com/laduke/zerotier-one_exporter/internal/promexporter"
	"github.com/laduke/zerotier-one_exporter/internal/ztoneservice"
)

// TODO make a config service, make a flag to collect zerotier addresses or not

func main() {
	var (
		webConfig     = webflag.AddFlags(kingpin.CommandLine)
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":19993").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		authtokenPath   = kingpin.Flag("zerotier.authtoken-path", "Path to the zerotier one authtoken.secret.").String()
	)


	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("zerotier_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	var token string
	if *authtokenPath != "" {
		b, err := ioutil.ReadFile(*authtokenPath)
		if err != nil { panic(err) }
		token = string(b)

	} else {
		token = authtoken.Guess()
	}
	oneService := ztone.New(token)

	reg := prometheus.NewPedanticRegistry()
	exporter := promexporter.New(oneService)
	reg.MustRegister(exporter)

	level.Info(logger).Log("msg", "Starting zerotier-one_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())


	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	http.Handle(*metricsPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>ZeroTier Exporter</title></head>
             <body>
             <h1>ZeroTier Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	srv := &http.Server{Addr: *listenAddress}

	if err := web.ListenAndServe(srv, *webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
