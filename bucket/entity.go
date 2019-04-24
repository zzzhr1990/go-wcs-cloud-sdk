package bucket

//StatResult response
type StatResult struct{
	Result bool `json:"result,omitempty"`
	Code int `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
	Fsize int64 `json:"fsize,omitempty"`
	Hash string `json:"hash,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	PutTime int64 `json:"putTime,omitempty"`
	ExpirationDate string `json:"expirationDate,omitempty"`
}
/*
{
    "result":  "true",
    "code":  "200",
    "name":  "<fileName string>",
    "message":  "< message string>",
    "fsize":     "<FileSize  int>",
    "hash":     "<FileETag  string>",
    "mimeType:  "<MimeType  string>",
    "putTime":    "<PutTime   int64>"
    "expirationDate":   "<ExpirationDate string>"
}
*/