package main

import (
	"encoding/json"

	"github.com/parnurzeal/gorequest"
	"github.com/prometheus/client_golang/prometheus"
)

type LightState struct {
	On               bool
	Brightness       int `json:"bri"`
	Hue              int
	Saturation       int `json:"sat"`
	Effect           string
	XY               []float64
	ColorTemperature int `json:"ct"`
	Alert            string
	ColorMode        string
	Mode             string
	Reachable        bool
}

type LightSoftwareUpdate struct {
	State       string
	LastInstall string
}

type Gamut []float32

type ColorTemperatureLimits struct {
	Min int
	Max int
}

type LightControl struct {
	MinDimLevel            int
	MaxLumen               int
	ColorGamutType         string
	ColorGamut             []Gamut
	ColorTemperatureLimits ColorTemperatureLimits `json:"ct"`
}

type LightStreaming struct {
	Renderer bool
	Proxy    bool
}

type LightCapabilities struct {
	Certified bool
	Control   LightControl
	Streaming LightStreaming
}

type LightStartup struct {
	Mode       string
	Configured bool
}

type LightConfig struct {
	Archetype string
	Function  string
	Direction string
	Startup   LightStartup
}

type LightAttributes struct {
	State            LightState
	SoftwareUpdate   LightSoftwareUpdate `json:"swupdate"`
	Type             string
	Name             string
	ModelId          string
	ManufacturerName string
	ProductName      string
	Capabilities     LightCapabilities
	Config           LightConfig
	UniqueId         string
	SoftwareVersion  string `json:"swversion"`
	SoftwareConfigId string `json:"swconfigid"`
	ProductId        string
}

type Lights map[int]LightAttributes

var (
	hueLightsStateDesc      = NewHueDesc("light", "state", "Whether the light is on or off.", "colormode")
	hueLightsReachableDesc  = NewHueDesc("light", "reachable", "Whether the light is reachable or not.")
	hueLightsBrightnessDesc = NewHueDesc("light", "brightness", "The brightness of the light.")
	hueLightsHDesc          = NewHueDesc("light", "color_hue", "The hue of the light's color.")
	hueLightsSDesc          = NewHueDesc("light", "color_saturation", "The saturation of the light's color.")
	hueLightsCTDesc         = NewHueDesc("light", "color_temperature", "The temperature of the light's color.")
	hueLightsXDesc          = NewHueDesc("light", "color_x", "The X value of the light's color.")
	hueLightsYDesc          = NewHueDesc("light", "color_y", "The Y value of the light's color.")
)

func lightsCollect(address string, userToken string, ch chan<- prometheus.Metric) error {
	request := gorequest.New()
	_, body, errs := request.Get(address + "/api/" + userToken + "/lights").
		End()
	if errs != nil {
		return errs[0]
	}
	var response Lights
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return err
	}
	for id, light := range response {
		if light.State.On {
			ch <- NewHueGauge(hueLightsStateDesc, 1.0, id, light.Name, light.State.ColorMode)
		} else {
			ch <- NewHueGauge(hueLightsStateDesc, 0.0, id, light.Name, light.State.ColorMode)
		}
		if light.State.Reachable {
			ch <- NewHueGauge(hueLightsReachableDesc, 1.0, id, light.Name)
		} else {
			ch <- NewHueGauge(hueLightsReachableDesc, 0.0, id, light.Name)
		}
		ch <- NewHueGauge(hueLightsBrightnessDesc, float64(light.State.Brightness), id, light.Name)
		switch light.State.ColorMode {
		case "ct":
			ch <- NewHueGauge(hueLightsCTDesc, float64(light.State.ColorTemperature), id, light.Name)
		case "xy":
			ch <- NewHueGauge(hueLightsXDesc, light.State.XY[0], id, light.Name)
			ch <- NewHueGauge(hueLightsYDesc, light.State.XY[1], id, light.Name)
		case "hs":
			ch <- NewHueGauge(hueLightsHDesc, float64(light.State.Hue), id, light.Name)
			ch <- NewHueGauge(hueLightsSDesc, float64(light.State.Saturation), id, light.Name)
		}
	}
	return nil
}
