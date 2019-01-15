package main

import (
	"encoding/json"

	"github.com/mickael-carl/hue_exporter/pkg/sensors"
	"github.com/mickael-carl/hue_exporter/pkg/util"
	"github.com/parnurzeal/gorequest"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	hueSensorStateDesc       = util.NewHueDesc("sensor", "state", "Whether the sensor is on or off.", "type")
	hueSensorReachableDesc   = util.NewHueDesc("sensor", "reachable", "Whether the sensor is reachable or not.", "type")
	hueSensorBatteryDesc     = util.NewHueDesc("sensor", "battery_percent", "Remaining battery in the sensor.", "type")
	hueSensorTemperatureDesc = util.NewHueDesc("sensor", "temperature_degrees", "Temperature measured by the sensor.", "type")
)

func sensorsCollect(address string, userToken string, ch chan<- prometheus.Metric) error {
	request := gorequest.New()
	_, body, errs := request.Get(address + "/api/" + userToken + "/sensors").
		End()
	if errs != nil {
		return errs[0]
	}
	var sensors sensors.Sensors
	if err := json.Unmarshal([]byte(body), &sensors); err != nil {
		return err
	}
	for id, sensor := range sensors {
		config := sensor.Config.(map[string]interface{})
		state := sensor.State.(map[string]interface{})
		switch sensor.Type {
		case "ZLLTemperature":
			if config["on"].(bool) {
				ch <- util.NewHueGauge(hueSensorStateDesc, 1.0, id, sensor.Name, "temperature")
			} else {
				ch <- util.NewHueGauge(hueSensorStateDesc, 0.0, id, sensor.Name, "temperature")
			}
			if config["reachable"].(bool) {
				ch <- util.NewHueGauge(hueSensorReachableDesc, 1.0, id, sensor.Name, "temperature")
			} else {
				ch <- util.NewHueGauge(hueSensorReachableDesc, 0.0, id, sensor.Name, "temperature")
			}
			ch <- util.NewHueGauge(hueSensorBatteryDesc, config["battery"].(float64), id, sensor.Name, "temperature")
			ch <- util.NewHueGauge(hueSensorTemperatureDesc, state["temperature"].(float64)/100, id, sensor.Name, "temperature")
		default:
		}
	}
	return nil
}
