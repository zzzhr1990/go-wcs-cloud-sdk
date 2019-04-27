package bucket

import (
	"errors"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/core"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
)

//Stat star
// 获取文件信息（stat）
// https://wcs.chinanetcenter.com/document/API/ResourceManage/stat
func (manager *Manager) Stat(bucket string, key string) (*StatResult, error) {
	if 0 == len(bucket) {
		err := errors.New("bucket is empty")
		return nil, err
	}
	if 0 == len(key) {
		err := errors.New("key is empty")
		return nil, err
	}

	url := manager.config.GetManageURLPrefix() + "/stat/" + utility.URLSafeEncodePair(bucket, key)
	request, err := utility.CreateGetRequest(url)
	if nil != err {
		return nil, err
	}
	res := &StatResult{}
	err = manager.httpManager.DoWithAuthRetry(request, manager.auth, res, 10)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Copy (copy)
// https://wcs.chinanetcenter.com/document/API/ResourceManage/copy
func (manager *Manager) Copy(src string, dst string) (*core.CommonResponse, error) {
	if 0 == len(src) {
		return nil, errors.New("src is empty")
	}
	if 0 == len(dst) {
		return nil, errors.New("dst is empty")
	}

	url := manager.config.GetManageURLPrefix() + "/copy/" + utility.URLSafeEncodeString(src) +
		"/" + utility.URLSafeEncodeString(dst)
	request, err := utility.CreatePostRequest(url)
	if nil != err {
		return nil, err
	}
	respEntity := &core.CommonResponse{}
	err = manager.httpManager.DoWithAuthRetry(request, manager.auth, respEntity, 10)
	if err != nil {
		return nil, err
	}
	return respEntity, nil
}

// Delete 删除文件（delete）
// https://wcs.chinanetcenter.com/document/API/ResourceManage/delete
func (manager *Manager) Delete(bucket string, key string) (*core.CommonResponse, error) {
	if 0 == len(bucket) {
		return nil, errors.New("bucket is empty")
	}
	if 0 == len(key) {
		return nil, errors.New("key is empty")
	}
	url := manager.config.GetManageURLPrefix() + "/delete/" + utility.URLSafeEncodePair(bucket, key)
	request, err := utility.CreatePostRequest(url)
	if nil != err {
		return nil, err
	}
	respEntity := &core.CommonResponse{}
	err = manager.httpManager.DoWithAuthRetry(request, manager.auth, respEntity, 10)
	if err != nil {
		return nil, err
	}
	return respEntity, nil
}
