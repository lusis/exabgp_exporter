package main

import (
	"bufio"
	"net/http"
	"os"

	"github.com/lusis/exabgp_exporter/pkg/exporter"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

var (
	exaBGPCLICommand = "exabgpcli"
	exaBGPCLIRoot    = "/etc/exabgp"
)

func main() {

	var (
		_             = kingpin.Command("stream", "run in stream mode (appropriate for embedding as an exabgp process)")
		shellCmd      = kingpin.Command("standalone", "run in standalone mode (calls exabgpcli on each scrape)").Default()
		exabgpcmd     = shellCmd.Flag("exabgp.cli.command", "exabgpcli command").Default(exaBGPCLICommand).String()
		exabgproot    = shellCmd.Flag("exabgp.root", "value of --root to be passed to exabgpcli").Default(exaBGPCLIRoot).String()
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9576").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("exabgp_exporter"))
	kingpin.HelpFlag.Short('h')

	switch kingpin.Parse() {
	case "standalone":
		log.Infof("starting exabgp_exporter %s in standalone mode using '%s --root %s'", version.Info(), *exabgpcmd, *exabgproot)
		log.Infof("build context: %s", version.BuildContext())
		e, err := exporter.NewStandaloneExporter(*exabgpcmd, *exabgproot)
		if err != nil {
			log.Fatal(err)
		}
		prometheus.MustRegister(e)
		prometheus.MustRegister(version.NewCollector("exabgp_exporter"))
	case "stream":
		log.Infof("starting exabgp_exporter %s in stream mode", version.Info())
		log.Infof("build context: %s", version.BuildContext())
		e, err := exporter.NewEmbeddedExporter()
		if err != nil {
			log.Fatal(err)
		}
		prometheus.MustRegister(e)
		prometheus.MustRegister(version.NewCollector("exabgp_exporter"))
		reader := bufio.NewReader(os.Stdin)
		e.Run(reader)
	}
	log.Infoln("Listening on", *listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
             <head><title>ExaBGP Exporter</title></head>
             <body>
             <h1>ExaBGP Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
