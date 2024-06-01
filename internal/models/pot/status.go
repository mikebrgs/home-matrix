package pot

import "time"

type PotStatus struct {
	TS          time.Time `json:"ts"`
	DeviceId    string    `json:"device_id"`
	Battery     float32   `json:"battery"`
	Memory      float32   `json:"memory"`
	CPU         float32   `json:"cpu"`
	Temperature float32   `json:"temperature"`
	Storage     float32   `json:"storage"`
}
