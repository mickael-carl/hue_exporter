package sensors

import "github.com/mickael-carl/hue_exporter/pkg/common"

type SensorState struct {
	Temperature *int
	Presence    *bool
	ButtonEvent *int
	Status      *int
	Flag        *bool
	LightLevel  *int
	Dark        *bool
	Daylight    *bool
	LastUpdated *string
}

type SensorConfig struct {
	On             bool
	Battery        *int
	Reachable      *bool
	Alert          *string
	LedIndication  *bool
	UserTest       *bool
	Pending        []string
	Sensitivity    *int
	SensitivityMax *int
	SunriseOffset  *int
	SunsetOffset   *int
}

type SensorCapabilities struct {
	Certified bool
}

type SensorAttributes struct {
	State            SensorState
	SoftwareUpdate   *common.SoftwareUpdate `json:"swupdate"`
	Config           SensorConfig
	Name             string
	Type             string
	ModelId          string
	ManufacturerName string
	ProductName      string
	SoftwareVersion  string `json:"swversion"`
	UniqueId         string
	Capabilities     SensorCapabilities
}

type Sensors map[int]SensorAttributes
