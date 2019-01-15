package lights

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
