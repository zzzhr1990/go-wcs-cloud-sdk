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
func NewConfig(useHTTPS bool, uploadHost string, manageHost string) ( *Config) {
	if 0 == len(uploadHost) {
		uploadHost = "apitestuser.up0.v1.wcsapi.com"
	}
	if 0 == len(manageHost) {
		manageHost = "apitestuser.mgr0.v1.wcsapi.com"
	}
	return &Config{useHTTPS, uploadHost, manageHost}
}

//NewDefaultConfig nc
func NewDefaultConfig() ( *Config) {
	return NewConfig(false, "", "")
}

//GetManageURLPrefix GMP
func (config *Config) GetManageURLPrefix() ( string) {
	if config.UseHTTPS {
		return "https://" + config.ManageHost
	} 
		return "http://" + config.ManageHost
	
}

//GetUploadURLPrefix GLP
func (config *Config) GetUploadURLPrefix() ( string) {
	if config.UseHTTPS {
		return "https://" + config.UploadHost
	}
	return	 "http://" + config.UploadHost
	
	
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