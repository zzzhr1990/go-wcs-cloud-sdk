package upload

/*
{
    "ctx":          "<Ctx           string>",
    "checksum":     "<Checksum      string>",
    "crc32":         "<Crc32         int64>",
    "offset":        "<Offset        int64>"
}
*/

// MakeBlockBputResult Result for make block
type MakeBlockBputResult struct {
	Code     int    `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
	Ctx      string `json:"ctx,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	Crc32    int64  `json:"crc32,omitempty"`
	Offset   string `json:"offset,omitempty"`
}

/*
{
    "hash":"<ContentHash>",
    "key":"<Key>"
}
*/

// MakeFileResult make
type MakeFileResult struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Hash    string `json:"hash,omitempty"`
	Key     string `json:"key,omitempty"`
}
