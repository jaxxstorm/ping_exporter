package main

import (
	"strings"
	"sync"

	mon "github.com/digineo/go-ping/monitor"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "ping_"

var (
	labelNames = []string{"target", "ip", "ip_version"}
	rttDesc    = prometheus.NewDesc(prefix+"rtt_ms", "Round trip time in millis (deprecated)", append(labelNames, "type"), nil)
	bestDesc   = prometheus.NewDesc(prefix+"rtt_best_ms", "Best round trip time in millis", labelNames, nil)
	worstDesc  = prometheus.NewDesc(prefix+"rtt_worst_ms", "Worst round trip time in millis", labelNames, nil)
	meanDesc   = prometheus.NewDesc(prefix+"rtt_mean_ms", "Mean round trip time in millis", labelNames, nil)
	medianDesc = prometheus.NewDesc(prefix+"rtt_median_ms", "Median round trip time in millis", labelNames, nil)
	stddevDesc = prometheus.NewDesc(prefix+"rtt_std_deviation_ms", "Standard deviation in millis", labelNames, nil)
	lossDesc   = prometheus.NewDesc(prefix+"loss_percent", "Packet loss in percent", labelNames, nil)
	mutex      = &sync.Mutex{}
)

type pingCollector struct {
	monitor *mon.Monitor
	metrics map[string]*mon.Metrics
}

func (p *pingCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rttDesc
	ch <- lossDesc
	ch <- bestDesc
	ch <- worstDesc
	ch <- meanDesc
	ch <- medianDesc
	ch <- stddevDesc
}

func (p *pingCollector) Collect(ch chan<- prometheus.Metric) {
	mutex.Lock()
	defer mutex.Unlock()

	metrics := p.monitor.Export()

	if len(metrics) > 0 {
		p.metrics = metrics
	}

	if p.metrics == nil || len(p.metrics) == 0 {
		return
	}

	for target, metrics := range p.metrics {
		l := strings.SplitN(target, " ", 3)

		if metrics.PacketsSent > metrics.PacketsLost {
			ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, float64(metrics.Best), append(l, "best")...)
			ch <- prometheus.MustNewConstMetric(bestDesc, prometheus.GaugeValue, float64(metrics.Best), l...)

			ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, float64(metrics.Worst), append(l, "worst")...)
			ch <- prometheus.MustNewConstMetric(worstDesc, prometheus.GaugeValue, float64(metrics.Worst), l...)

			ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, float64(metrics.Mean), append(l, "mean")...)
			ch <- prometheus.MustNewConstMetric(meanDesc, prometheus.GaugeValue, float64(metrics.Mean), l...)

			ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, float64(metrics.Median), append(l, "median")...)
			ch <- prometheus.MustNewConstMetric(medianDesc, prometheus.GaugeValue, float64(metrics.Median), l...)

			ch <- prometheus.MustNewConstMetric(rttDesc, prometheus.GaugeValue, float64(metrics.StdDev), append(l, "std_dev")...)
			ch <- prometheus.MustNewConstMetric(stddevDesc, prometheus.GaugeValue, float64(metrics.StdDev), l...)
		}

		loss := float64(metrics.PacketsLost) / float64(metrics.PacketsSent)
		ch <- prometheus.MustNewConstMetric(lossDesc, prometheus.GaugeValue, loss, l...)
	}
}
