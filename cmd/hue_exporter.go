package main

import (
	"net/http"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

var (
	hueScrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName("hue", "", "scrape_success"),
		"Whether scraping the Hue devices was successful.",
		nil, nil)
	app           = kingpin.New("hue_exporter", "A Prometheus exporter for the Philips Hue IoT system.")
	bridgeAddress = app.Flag("bridge", "Bridge HTTP address.").Required().String()
	userToken     = app.Flag("token", "The user token used to authenticate requests against the bridge.").Required().String()
	listenAddress = app.Flag("listen-address", "Address on which to expose metrics.").Default(":9535").String()
	metricsPath   = app.Flag("metrics-path", "Path under which to expose metrics.").Default("/metrics").String()
)

type HueExporter struct {
	address   string
	userToken string
}

func (e *HueExporter) Collect(ch chan<- prometheus.Metric) {
	err := lightsCollect(e.address, e.userToken, ch)
	if err != nil {
		log.Error("Failed to gather info on Hue system: ", err)
		ch <- prometheus.MustNewConstMetric(
			hueScrapeSuccessDesc, prometheus.GaugeValue, 0.0)
	}
	err = sensorsCollect(e.address, e.userToken, ch)
	if err != nil {
		log.Error("Failed to gather info on Hue system: ", err)
		ch <- prometheus.MustNewConstMetric(
			hueScrapeSuccessDesc, prometheus.GaugeValue, 0.0)
	}
	err = groupsCollect(e.address, e.userToken, ch)
	if err != nil {
		log.Error("Failed to gather info on Hue system: ", err)
		ch <- prometheus.MustNewConstMetric(
			hueScrapeSuccessDesc, prometheus.GaugeValue, 0.0)
	} else {
		ch <- prometheus.MustNewConstMetric(
			hueScrapeSuccessDesc, prometheus.GaugeValue, 1.0)
	}
}

func (e *HueExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- hueScrapeSuccessDesc
}

func newHueExporter(address string, userToken string) (*HueExporter, error) {
	return &HueExporter{
		address:   address,
		userToken: userToken,
	}, nil
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	exporter, err := newHueExporter(*bridgeAddress, *userToken)
	if err != nil {
		log.Error(err)
	}
	log.Infoln("Starting hue_exporter")
	log.Infoln("Exporting metrics from Hue bridge located at", exporter.address)
	prometheus.MustRegister(exporter)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Hue Exporter</title></head>
			<body>
			<h1>Hue Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
