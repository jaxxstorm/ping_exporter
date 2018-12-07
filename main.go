package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/czerwonk/ping_exporter/config"
	"github.com/digineo/go-ping"
	mon "github.com/digineo/go-ping/monitor"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

var version string = "snapshot"
var commit string = "unknown"
var date string = "now"

var (
	showVersion   = kingpin.Flag("version", "Print version information").Default().Bool()
	listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface").Default(":9427").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics").Default("/metrics").String()
	configFile    = kingpin.Flag("config.path", "Path to config file").Default("").String()
	pingInterval  = kingpin.Flag("ping.interval", "Interval for ICMP echo requests").Default("5s").Duration()
	pingTimeout   = kingpin.Flag("ping.timeout", "Timeout for ICMP echo request").Default("4s").Duration()
	historySize   = kingpin.Flag("ping.history-size", "Number of results to remember per target").Default("10").Int()
	dnsRefresh    = kingpin.Flag("dns.refresh", "Interval for refreshing DNS records and updating targets accordingly (0 if disabled)").Default("1m").Duration()
	dnsNameServer = kingpin.Flag("dns.nameserver", "DNS server used to resolve hostname of targets").Default("").String()
	logLevel      = kingpin.Flag("log.level", "Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]").Default("info").String()
	targets       = kingpin.Arg("targets", "A list of targets to ping").Strings()
)

func init() {
	kingpin.Parse()
}

func main() {

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	if *historySize < 1 {
		fmt.Println("ping.history-size must be greater than 0")
		os.Exit(0)
	}

	err := log.Logger.SetLevel(log.Base(), *logLevel)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	if len(cfg.Targets) == 0 {
		kingpin.FatalUsage("No targets specified")
	}

	m, err := startMonitor(cfg)
	if err != nil {
		log.Errorln(err)
		os.Exit(2)
	}

	startServer(m)
}

func printVersion() {
	fmt.Println("ping-exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Build Date: %s\n", date)
	fmt.Println("Author(s): Philip Berndroth, Daniel Czerwonk")
	fmt.Println("Metric exporter for go-icmp")
}

func startMonitor(cfg *config.Config) (*mon.Monitor, error) {
	pinger, err := ping.New("0.0.0.0", "::")
	if err != nil {
		return nil, err
	}

	monitor := mon.New(pinger, *pingInterval, *pingTimeout)
	monitor.HistorySize = *historySize
	targets := make([]*target, len(cfg.Targets))
	for i, host := range cfg.Targets {
		t := &target{
			host:      host,
			addresses: make([]net.IP, 0),
			delay:     time.Duration(10*i) * time.Millisecond,
			dns:       *dnsNameServer,
		}
		targets[i] = t

		err := t.addOrUpdateMonitor(monitor)
		if err != nil {
			log.Errorln(err)
		}
	}

	go startDNSAutoRefresh(targets, monitor)

	return monitor, nil
}

func startDNSAutoRefresh(targets []*target, monitor *mon.Monitor) {
	if *dnsRefresh == 0 {
		return
	}

	for {
		select {
		case <-time.After(*dnsRefresh):
			refreshDNS(targets, monitor)
		}
	}
}

func refreshDNS(targets []*target, monitor *mon.Monitor) {
	for _, t := range targets {
		log.Infoln("refreshing DNS")

		go func(ta *target) {
			err := ta.addOrUpdateMonitor(monitor)
			if err != nil {
				log.Errorf("could refresh dns: %v", err)
			}
		}(t)
	}
}

func startServer(monitor *mon.Monitor) {
	log.Infof("Starting ping exporter (Version: %s)", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>ping Exporter (Version ` + version + `)</title></head>
			<body>
			<h1>ping Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			<h2>More information:</h2>
			<p><a href="https://github.com/czerwonk/ping_exporter">github.com/czerwonk/ping_exporter</a></p>
			</body>
			</html>`))
	})

	reg := prometheus.NewRegistry()
	reg.MustRegister(&pingCollector{monitor: monitor})
	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorLog:      log.NewErrorLogger(),
		ErrorHandling: promhttp.ContinueOnError})
	http.HandleFunc("/metrics", h.ServeHTTP)

	log.Infof("Listening for %s on %s", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func loadConfig() (*config.Config, error) {
	if *configFile == "" {
		return &config.Config{Targets: *targets}, nil
	}

	b, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	return config.FromYAML(bytes.NewReader(b))
}
