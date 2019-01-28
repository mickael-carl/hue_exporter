package main

import (
	"crypto/tls"
	"encoding/json"

	"github.com/mickael-carl/hue_exporter/pkg/groups"
	"github.com/mickael-carl/hue_exporter/pkg/util"
	"github.com/parnurzeal/gorequest"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	hueGroupsLightsNumberDesc  = util.NewHueDesc("group", "lights", "The number of lights belonging to the group.", "type")
	hueGroupsSensorsNumberDesc = util.NewHueDesc("group", "sensors", "The number of sensors belonging to the group.", "type")
)

func groupsCollect(address string, tlsConfig *tls.Config, userToken string, ch chan<- prometheus.Metric) error {
	request := gorequest.New()
	_, body, errs := request.Get(address + "/api/" + userToken + "/groups").
		TLSClientConfig(tlsConfig).
		End()
	if errs != nil {
		return errs[0]
	}
	var response groups.Groups
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return err
	}
	for id, group := range response {
		ch <- util.NewHueGauge(hueGroupsLightsNumberDesc, float64(len(group.Lights)), id, group.Name, group.GroupType)
		ch <- util.NewHueGauge(hueGroupsSensorsNumberDesc, float64(len(group.Sensors)), id, group.Name, group.GroupType)
	}
	return nil
}
