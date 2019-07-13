package entity

// UploadPolicy upp
type UploadPolicy struct {
	Scope               string `json:"scope,omitempty"`    //bucket
	Deadline            string `json:"deadline,omitempty"` //upload deadline
	SaveKey             string `json:"saveKey,omitempty"`
	ReturnURL           string `json:"returnUrl,omitempty"`
	ReturnBody          string `json:"returnBody,omitempty"`
	Overwrite           int    `json:"overwrite"`
	FsizeLimit          int64  `json:"fsizeLimit,omitempty"`
	CallbackURL         string `json:"callbackUrl,omitempty"`
	CallbackBody        string `json:"callbackBody,omitempty"`
	PersistentOps       string `json:"persistentOps,omitempty"`
	PersistentNotifyURL string `json:"persistentNotifyUrl,omitempty"`
	// ContentDetect string
	// DetectNotifyURL
	// detectNotifyRule
	Separate string `json:"separate,omitempty"`
}
