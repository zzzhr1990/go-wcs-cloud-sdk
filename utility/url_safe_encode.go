package utility

import (
	"bytes"
	"encoding/base64"
	"net/url"
	"sort"
)

//URLSafeEncode URLSafeEncode
func URLSafeEncode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

//URLSafeEncodeString c
func URLSafeEncodeString(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

//URLSafeEncodePair /d
func URLSafeEncodePair(bucket string, key string) string {
	return base64.URLEncoding.EncodeToString([]byte(bucket + ":" + key))
}

//MakeQuery m
func MakeQuery(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}
