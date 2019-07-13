package core

// ListResponse response
type ListResponse struct {
	Code           int      `json:"code,omitempty"`
	Message        string   `json:"message,omitempty"`
	Marker         string   `json:"marker,omitempty"`
	CommonPrefixes []string `json:"commonPrefixes,omitempty"`

	Items []struct {
		Key            string `json:"key,omitempty"`
		PutTime        int64  `json:"putTime,omitempty"`
		Hash           string `json:"hash,omitempty"`
		FSize          string `json:"fsize,omitempty"`
		MimeType       string `json:"mimeType,omitempty"`
		ExpirationDate string `json:"expirationDate,omitempty"`
	} `json:"items"`
}
