package bucket

import (
	"errors"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
)

//Stat star
// 获取文件信息（stat）
// https://wcs.chinanetcenter.com/document/API/ResourceManage/stat
func (manager *Manager) Stat(bucket string, key string) ( *StatResult,  error) {
	if 0 == len(bucket) {
		err := errors.New("bucket is empty")
		return nil,err
	}
	if 0 == len(key) {
		err := errors.New("key is empty")
		return nil,err
	}

	url := manager.config.GetManageURLPrefix() + "/stat/" + utility.URLSafeEncodePair(bucket, key)
	request, err := utility.CreateGetRequest(url)
	if nil != err {
		return nil,err
	}
	res :=&StatResult{}
	err = manager.httpManager.DoWithAuthRetry(request, manager.auth, res, 10)
	if err != nil{
		return nil, err
	}
	return res, nil
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