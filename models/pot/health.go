package pot

type PotHealth struct {
	TS          int64   `json:"ts"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	Pressure    float32 `json:"pressure"`
	Moisture    float32 `json:"moisture"`
	Light       float32 `json:"light"`
}
