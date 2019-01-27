package main

import (
	"encoding/json"
	"math"

	"github.com/mickael-carl/hue_exporter/pkg/sensors"
	"github.com/mickael-carl/hue_exporter/pkg/util"
	"github.com/parnurzeal/gorequest"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	hueSensorStateDesc         = util.NewHueDesc("sensor", "state", "Whether the sensor is on or off.", "type")
	hueSensorReachableDesc     = util.NewHueDesc("sensor", "reachable", "Whether the sensor is reachable or not.", "type")
	hueSensorBatteryDesc       = util.NewHueDesc("sensor", "battery_percent", "Remaining battery in the sensor.", "type")
	hueSensorTemperatureDesc   = util.NewHueDesc("sensor", "temperature_degrees", "Temperature measured by the sensor.", "type")
	hueSensorPresenceDesc      = util.NewHueDesc("sensor", "presence", "Whether presence is detected by the sensor.", "type")
	hueSensorLightLevelDesc    = util.NewHueDesc("sensor", "light_level_lumens", "Light level measured by the sensor.", "type")
	hueSensorDaylightDesc      = util.NewHueDesc("sensor", "daylight", "Weather current time is between sunrise and sunset at current location.", "type")
	hueSensorButtonPressedDesc = util.NewHueDesc("sensor", "button_pressed", "Code of the last switch event of the sensor.", "type")
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
		if sensor.Config.On {
			ch <- util.NewHueGauge(hueSensorStateDesc, 1.0, id, sensor.Name, sensor.Type)
		} else {
			ch <- util.NewHueGauge(hueSensorStateDesc, 0.0, id, sensor.Name, sensor.Type)
		}
		if sensor.Config.Reachable != nil {
			if *sensor.Config.Reachable {
				ch <- util.NewHueGauge(hueSensorReachableDesc, 1.0, id, sensor.Name, sensor.Type)
			} else {
				ch <- util.NewHueGauge(hueSensorReachableDesc, 0.0, id, sensor.Name, sensor.Type)
			}
		}
		if sensor.Config.Battery != nil {
			ch <- util.NewHueGauge(hueSensorBatteryDesc, float64(*sensor.Config.Battery), id, sensor.Name, sensor.Type)
		}
		if sensor.State.Temperature != nil {
			ch <- util.NewHueGauge(hueSensorTemperatureDesc, float64(*sensor.State.Temperature)/100, id, sensor.Name, sensor.Type)
		}
		if sensor.State.Presence != nil {
			if *sensor.State.Presence {
				ch <- util.NewHueGauge(hueSensorPresenceDesc, 1.0, id, sensor.Name, sensor.Type)
			} else {
				ch <- util.NewHueGauge(hueSensorPresenceDesc, 0.0, id, sensor.Name, sensor.Type)
			}
		}
		if sensor.State.LightLevel != nil {
			ch <- util.NewHueGauge(hueSensorLightLevelDesc, math.Pow(10, (float64(*sensor.State.LightLevel)-1)/10000), id, sensor.Name, sensor.Type)
		}
		if sensor.State.Daylight != nil {
			if *sensor.State.Daylight {
				ch <- util.NewHueGauge(hueSensorDaylightDesc, 1.0, id, sensor.Name, sensor.Type)
			} else {
				ch <- util.NewHueGauge(hueSensorDaylightDesc, 0.0, id, sensor.Name, sensor.Type)
			}
		}
		if sensor.State.ButtonEvent != nil {
			ch <- util.NewHueGauge(hueSensorButtonPressedDesc, float64(*sensor.State.ButtonEvent), id, sensor.Name, sensor.Type)
		}

	}
	return nil
}
