package pot

import "time"

type PotHealth struct {
	TS          time.Time `json:"ts"`
	DeviceId    string    `json:"device_id"`
	Temperature float32   `json:"temperature"`
	Humidity    float32   `json:"humidity"`
	Pressure    float32   `json:"pressure"`
	Moisture    float32   `json:"moisture"`
	Light       float32   `json:"light"`
}
