package core

import (
	"github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
)

//Config ce
type Config struct {
	UseHTTPS   bool
	UploadHost string
	ManageHost string
}

//NewConfig nc
func NewConfig(useHTTP bool, uploadHost string, manageHost string) (config *Config) {
	if 0 == len(uploadHost) {
		uploadHost = "apitestuser.up0.v1.wcsapi.com"
	}
	if 0 == len(manageHost) {
		manageHost = "apitestuser.mgr0.v1.wcsapi.com"
	}
	return &Config{useHTTP, uploadHost, manageHost}
}

//NewDefaultConfig nc
func NewDefaultConfig() (config *Config) {
	return NewConfig(false, "", "")
}

//GetManageURLPrefix GMP
func (config *Config) GetManageURLPrefix() (urlPrefix string) {
	if config.UseHTTPS {
		urlPrefix = "https://" + config.ManageHost
	} else {
		urlPrefix = "http://" + config.ManageHost
	}
	return
}

//GetUploadURLPrefix GLP
func (config *Config) GetUploadURLPrefix() (urlPrefix string) {
	if config.UseHTTPS {
		urlPrefix = "https://" + config.UploadHost
	} else {
		urlPrefix = "http://" + config.UploadHost
	}
	return
}

const (
	//Version c
	Version        = "1.0.0.1" //Version vd
	//BlockSize bs
	BlockSize int = 4 * 1024 * 1024//BlockSize bs 
)

func init() {
	utility.SetUserAgent("WCS-GO-SDK-" + Version)
}