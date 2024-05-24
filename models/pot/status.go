package pot

type PotStatus struct {
	TS          int64   `json:"ts"`
	Battery     float32 `json:"battery"`
	Memory      float32 `json:"memory"`
	CPU         float32 `json:"cpu"`
	Temperature float32 `json:"temperature"`
	Storage     float32 `json:"storage"`
}
