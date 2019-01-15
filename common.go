package main

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func NewHueDesc(deviceType string, metricName string, description string, otherLabels ...string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName("hue", deviceType, metricName),
		description,
		append([]string{"id", "name"}, otherLabels...),
		nil)
}

func NewHueGauge(desc *prometheus.Desc, value float64, id int, name string, otherLabels ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value,
		append([]string{strconv.Itoa(id), name}, otherLabels...)...)
}
