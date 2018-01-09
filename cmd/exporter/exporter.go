package main

import (
	"fmt"
	"net/http"
	"os"

	flag "github.com/bborbe/flagenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/seibert-media/jenkins_exporter/pkg/exporter"

	_ "net/http/pprof"
)

var (
	showVersion   = flag.Bool("version", false, "Print version information")
	listenAddress = flag.String("web-listen-address", ":9103", "Address to listen on for web interface and telemetry")
	metricsPath   = flag.String("web-telemetry-path", "/metrics", "Path to expose metrics of the exporter")
	address       = flag.String("jenkins-address", "", "Address where to reach Jenkins")
	username      = flag.String("jenkins-username", "", "Username to authenticate on Jenkins")
	password      = flag.String("jenkins-password", "", "Password to authenticate on Jenkins")
	debug         = flag.Bool("debug", false, "debug logging")
)

// init registers the collector version.
func init() {
	prometheus.MustRegister(version.NewCollector("jenkins_exporter"))
}

// main simply initializes this tool.
func main() {
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("jenkins_exporter"))
		os.Exit(0)
	}

	if *address == "" {
		fmt.Fprintln(os.Stderr, "Please provide a address for Jenkins")
		os.Exit(1)
	}

	log.Infoln("Starting Jenkins exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	e := exporter.NewExporter(*address, *username, *password)

	prometheus.MustRegister(e)
	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.Unregister(prometheus.NewProcessCollector(os.Getpid(), ""))

	http.Handle(*metricsPath, promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Infof("Listening on %s", *listenAddress)

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatal(err)
	}
}
