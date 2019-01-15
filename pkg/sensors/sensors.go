package sensors

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

type Sensors map[int]SensorAttributes
