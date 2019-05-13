package utility

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	// "github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
)

// CommonRequest for create request
type CommonRequest struct {
	uri        string
	method     string
	header     map[string]string
	token      string
	auth       *Auth
	data       []byte
	stringBody string
	timeOut    time.Duration
}

// CreateCommonRequest cra
func CreateCommonRequest(method string, uri string) (*CommonRequest, error) {
	return &CommonRequest{
		uri:    uri,
		method: method,
		header: make(map[string]string),
	}, nil
}

// CreateRequest ccr
func (req *CommonRequest) CreateRequest() (*http.Request, error) {
	//
	var reader io.Reader
	if req.data != nil && len(req.data) > 0 {
		reader = bytes.NewReader(req.data)
	} else {
		if len(req.stringBody) > 0 {
			reader = strings.NewReader(req.stringBody)
		}
	}
	request, err := http.NewRequest(req.method, req.uri, reader)

	if err != nil {
		log.Errorf("create request failed, %v:%v", req.method, req.uri)
		return nil, err
	}

	for k, v := range req.header {
		request.Header.Set(k, v)
	}

	if _, ok := request.Header["User-Agent"]; !ok {
		request.Header.Set("User-Agent", "wcs-common-0.0.1")
	}

	if len(req.token) > 0 {
		request.Header.Set("Authorization", req.token)
	} else {
		if req.auth != nil {
			var token string
			token, err = req.auth.SignRequest(request)
			if nil != err {
				return nil, err
			}
			request.Header.Set("Authorization", token)
		}
	}

	return request, nil
}

// AddToken add token
func (req *CommonRequest) AddToken(token string) {
	req.token = token
}

// AddAuth add token
func (req *CommonRequest) AddAuth(auth *Auth) {
	req.auth = auth
}

// AddHeader add header
func (req *CommonRequest) AddHeader(header string, value string) {
	req.header[header] = value
}

// AddData data
func (req *CommonRequest) AddData(data []byte) {
	req.data = data
}

// AddStringBody data
func (req *CommonRequest) AddStringBody(data string) {
	req.stringBody = data
}

// SetTimeout data
func (req *CommonRequest) SetTimeout(data time.Duration) {
	req.timeOut = data
}

// GetTimeout get request timeout
func (req *CommonRequest) GetTimeout() time.Duration {
	return req.timeOut
}
