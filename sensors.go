package main

import (
	"encoding/json"

	"github.com/parnurzeal/gorequest"
	"github.com/prometheus/client_golang/prometheus"
)

type TempSensorState struct {
	Temperature int
	LastUpdated string
}

type TempSensorConfig struct {
	On            bool
	Battery       int
	Reachable     bool
	Alert         string
	LedIndication bool
	UserTest      bool
	Pending       interface{}
}

type SensorAttributes struct {
	State            interface{}
	SwUpdate         interface{}
	Config           interface{}
	Name             string
	Type             string
	ModelId          string
	ManufacturerName string
	ProductName      string
	SoftwareVersion  string `json:"swversion"`
	UniqueId         string
	Capabilities     interface{}
}

//type DaylightSensorAttributes struct {
//	State DaylightSensorState
//	Config DaylightSensorConfig
//	SensorsCommonAttributes
//}

type Sensors map[int]SensorAttributes

var (
	hueSensorStateDesc       = NewHueDesc("sensor", "state", "Whether the sensor is on or off.", "type")
	hueSensorReachableDesc   = NewHueDesc("sensor", "reachable", "Whether the sensor is reachable or not.", "type")
	hueSensorBatteryDesc     = NewHueDesc("sensor", "battery_percent", "Remaining battery in the sensor.", "type")
	hueSensorTemperatureDesc = NewHueDesc("sensor", "temperature_degrees", "Temperature measured by the sensor.", "type")
)

func sensorsCollect(address string, userToken string, ch chan<- prometheus.Metric) error {
	request := gorequest.New()
	_, body, errs := request.Get(address + "/api/" + userToken + "/sensors").
		End()
	if errs != nil {
		return errs[0]
	}
	var sensors Sensors
	if err := json.Unmarshal([]byte(body), &sensors); err != nil {
		return err
	}
	for id, sensor := range sensors {
		config := sensor.Config.(map[string]interface{})
		state := sensor.State.(map[string]interface{})
		switch sensor.Type {
		case "ZLLTemperature":
			if config["on"].(bool) {
				ch <- NewHueGauge(hueSensorStateDesc, 1.0, id, sensor.Name, "temperature")
			} else {
				ch <- NewHueGauge(hueSensorStateDesc, 0.0, id, sensor.Name, "temperature")
			}
			if config["reachable"].(bool) {
				ch <- NewHueGauge(hueSensorReachableDesc, 1.0, id, sensor.Name, "temperature")
			} else {
				ch <- NewHueGauge(hueSensorReachableDesc, 0.0, id, sensor.Name, "temperature")
			}
			ch <- NewHueGauge(hueSensorBatteryDesc, config["battery"].(float64), id, sensor.Name, "temperature")
			ch <- NewHueGauge(hueSensorTemperatureDesc, state["temperature"].(float64)/100, id, sensor.Name, "temperature")
		default:
		}
	}
	return nil
}
