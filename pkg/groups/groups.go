package groups

type GroupAttributes struct {
	Action    GroupAction
	State     GroupState
	Lights    []string
	Sensors   []string
	Name      string
	GroupType string `json:"type"`
	ModelId   string
	UniqueId  string
	Class     string
	Recycle   bool
}

type GroupAction struct {
	On               bool
	Brightness       int `json:"bri"`
	Hue              int
	Saturation       int `json:"sat"`
	Effect           string
	XY               []float64
	ColorTemperature int `json:"ct"`
	Alert            string
	ColorMode        string
}

type GroupState struct {
	AllOn bool
	AnyOn bool
}

type Groups map[int]GroupAttributes
