package core

// FopsResponse response
type FopsResponse struct {
	Code         int    `json:"code,omitempty"`
	Message      string `json:"message,omitempty"`
	PersistentID string `json:"persistentId"`
}

// DetectResponse response
type DetectResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Results []struct {
		Image      string `json:"image,omitempty"`
		ErrMessage string `json:"errMessage,omitempty"`
		PornDetect struct {
			Label  int32   `json:"label,omitempty"`
			Rate   float64 `json:"rate,omitempty"`
			Review bool    `json:"review,omitempty"`
		} `json:"pornDetect,omitempty"`
		TerrorDetect struct {
			Label  int32   `json:"label,omitempty"`
			Rate   float64 `json:"rate,omitempty"`
			Review bool    `json:"review,omitempty"`
		} `json:"terrorDetect,omitempty"`
		PoliticalDetect struct {
			Label   int32 `json:"label,omitempty"`
			Persons []struct {
				Name   string  `json:"name,omitempty"`
				Rate   float64 `json:"rate,omitempty"`
				Review bool    `json:"review,omitempty"`
			} `json:"persons,omitempty"`
		} `json:"politicalDetect,omitempty"`
	} `json:"results"`
}
