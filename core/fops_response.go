package core

// FopsResponse response
type FopsResponse struct {
	Code         int    `json:"code,omitempty"`
	Message      string `json:"message,omitempty"`
	PersistentID string `json:"persistentId"`
}
