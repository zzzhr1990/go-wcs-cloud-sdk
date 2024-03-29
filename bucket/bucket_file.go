package bucket

import (
	"errors"
	"strconv"

	"fmt"

	"github.com/zzzhr1990/go-wcs-cloud-sdk/core"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/wcserror"
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

// Move (copy)
// https://wcs.chinanetcenter.com/document/API/ResourceManage/copy
func (manager *Manager) Move(src string, dst string) (*core.CommonResponse, error) {
	return manager.MoveWithRetry(src, dst, 10)
}

// MoveWithRetry (copy)
// https://wcs.chinanetcenter.com/document/API/ResourceManage/copy
func (manager *Manager) MoveWithRetry(src string, dst string, retry int) (*core.CommonResponse, error) {
	if 0 == len(src) {
		return nil, errors.New("src is empty")
	}
	if 0 == len(dst) {
		return nil, errors.New("dst is empty")
	}

	url := manager.config.GetManageURLPrefix() + "/move/" + utility.URLSafeEncodeString(src) +
		"/" + utility.URLSafeEncodeString(dst)
	request, err := utility.CreatePostRequest(url)
	if nil != err {
		return nil, err
	}
	respEntity := &core.CommonResponse{}
	err = manager.httpManager.DoWithAuthRetry(request, manager.auth, respEntity, retry)
	if err != nil {
		return nil, err
	}
	return respEntity, nil
}

// List 删除文件（delete）
// https://wcs.chinanetcenter.com/document/API/ResourceManage/delete
func (manager *Manager) List(bucket string, limit int64, prefix string, mode int, marker string, retry int) (*core.ListResponse, error) {
	// return manager.DeleteWithRetry(bucket, key, 10)

	var queryStr string
	queryStr += "bucket=" + bucket + "&"
	if limit >= 1 && limit <= 1000 {
		queryStr += "limit=" + strconv.FormatInt(limit, 10) + "&"
	}

	if len(prefix) > 0 {
		queryStr += "prefix=" + utility.URLSafeEncodeString(prefix) + "&"
	}

	if 0 == mode || 1 == mode {
		queryStr += "mode=" + strconv.Itoa(mode) + "&"
	}

	if len(marker) > 0 {
		queryStr += "marker=" + marker
	}
	url := manager.config.GetManageURLPrefix() + "/list?" + queryStr
	// &mode=<mode>
	request, err := utility.CreateGetRequest(url)
	// fmt.Println(url)
	if nil != err {
		return nil, err
	}
	respEntity := &core.ListResponse{}
	err = manager.httpManager.DoWithAuthRetry(request, manager.auth, respEntity, retry)
	if err != nil {
		return nil, err
	}
	return respEntity, nil
}

// DeleteAll 删除文件（delete）
// https://wcs.chinanetcenter.com/document/API/ResourceManage/delete
func (manager *Manager) DeleteAll(bucket string, prefix string) (bool, error) {
	// return manager.DeleteWithRetry(bucket, key, 10)

	marker := ""
	for {
		lr, err := manager.List(bucket, 1000, prefix, 0, marker, 10)
		if err != nil {
			return false, err
		}
		if len(lr.Items) < 1 {
			return true, nil
		}
		for _, itm := range lr.Items {
			key := itm.Key
			// bucket := itm.
			go func(k string, v string) {
				_, err = manager.Delete(bucket, key)
				if err != nil && err != wcserror.ErrFileNotFound {
					fmt.Println(err)
				}
			}(bucket, key)
			// _, err = manager.Delete(bucket, key)
			// if err != nil && err != wcserror.ErrFileNotFound {
			// 	return false, err
			//}
		}
		marker = lr.Marker
	}
}

// Delete 删除文件（delete）
// https://wcs.chinanetcenter.com/document/API/ResourceManage/delete
func (manager *Manager) Delete(bucket string, key string) (*core.CommonResponse, error) {
	return manager.DeleteWithRetry(bucket, key, 10)
}

// DeleteWithRetry retry
func (manager *Manager) DeleteWithRetry(bucket string, key string, retry int) (*core.CommonResponse, error) {
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
	err = manager.httpManager.DoWithAuthRetry(request, manager.auth, respEntity, retry)
	if err != nil {
		return nil, err
	}
	return respEntity, nil
}
