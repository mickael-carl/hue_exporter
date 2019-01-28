package main

import (
	"crypto/tls"
	"encoding/json"

	"github.com/mickael-carl/hue_exporter/pkg/lights"
	"github.com/mickael-carl/hue_exporter/pkg/util"
	"github.com/parnurzeal/gorequest"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	hueLightsStateDesc      = util.NewHueDesc("light", "state", "Whether the light is on or off.", "colormode")
	hueLightsReachableDesc  = util.NewHueDesc("light", "reachable", "Whether the light is reachable or not.")
	hueLightsBrightnessDesc = util.NewHueDesc("light", "brightness", "The brightness of the light.")
	hueLightsHDesc          = util.NewHueDesc("light", "color_hue", "The hue of the light's color.")
	hueLightsSDesc          = util.NewHueDesc("light", "color_saturation", "The saturation of the light's color.")
	hueLightsCTDesc         = util.NewHueDesc("light", "color_temperature", "The temperature of the light's color.")
	hueLightsXDesc          = util.NewHueDesc("light", "color_x", "The X value of the light's color.")
	hueLightsYDesc          = util.NewHueDesc("light", "color_y", "The Y value of the light's color.")
)

func lightsCollect(address string, tlsConfig *tls.Config, userToken string, ch chan<- prometheus.Metric) error {
	request := gorequest.New()
	_, body, errs := request.Get(address + "/api/" + userToken + "/lights").
		TLSClientConfig(tlsConfig).
		End()
	if errs != nil {
		return errs[0]
	}
	var response lights.Lights
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return err
	}
	for id, light := range response {
		if light.State.On {
			ch <- util.NewHueGauge(hueLightsStateDesc, 1.0, id, light.Name, light.State.ColorMode)
		} else {
			ch <- util.NewHueGauge(hueLightsStateDesc, 0.0, id, light.Name, light.State.ColorMode)
		}
		if light.State.Reachable {
			ch <- util.NewHueGauge(hueLightsReachableDesc, 1.0, id, light.Name)
		} else {
			ch <- util.NewHueGauge(hueLightsReachableDesc, 0.0, id, light.Name)
		}
		ch <- util.NewHueGauge(hueLightsBrightnessDesc, float64(light.State.Brightness), id, light.Name)
		switch light.State.ColorMode {
		case "ct":
			ch <- util.NewHueGauge(hueLightsCTDesc, float64(light.State.ColorTemperature), id, light.Name)
		case "xy":
			ch <- util.NewHueGauge(hueLightsXDesc, light.State.XY[0], id, light.Name)
			ch <- util.NewHueGauge(hueLightsYDesc, light.State.XY[1], id, light.Name)
		case "hs":
			ch <- util.NewHueGauge(hueLightsHDesc, float64(light.State.Hue), id, light.Name)
			ch <- util.NewHueGauge(hueLightsSDesc, float64(light.State.Saturation), id, light.Name)
		}
	}
	return nil
}
