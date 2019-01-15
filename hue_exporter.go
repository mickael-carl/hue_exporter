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
	app = kingpin.New("hue_exporter", "A Prometheus exporter for the Philips Hue IoT system.")
	bridgeAddress = app.Flag("bridge", "Bridge HTTP address.").Required().String()
	userToken = app.Flag("token", "The user token used to authenticate requests against the bridge.").Required().String()
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
	prometheus.MustRegister(exporter)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
