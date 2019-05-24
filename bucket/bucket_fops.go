package bucket

import (
	"errors"
	"strings"

	"github.com/zzzhr1990/go-wcs-cloud-sdk/core"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
)

// Fops dp fops
func (manager *Manager) Fops(query string) (*core.FopsResponse, error) {
	if 0 == len(query) {
		return nil, errors.New("query is empty")
	}

	url := manager.config.GetManageURLPrefix() + "/fops"
	request, err := utility.CreatePostRequest(url)
	if nil != err {
		return nil, err
	}
	request.AddStringBody(query)
	// log.Infof("call fops %v", query)
	respEntity := &core.FopsResponse{}
	err = manager.httpManager.DoWithAuthRetry(request, manager.auth, respEntity, 10)
	if err != nil {
		return nil, err
	}
	return respEntity, nil
}

// FopsMap dp fops
func (manager *Manager) FopsMap(query map[string]string) (*core.FopsResponse, error) {
	if query == nil || 0 == len(query) {
		return nil, errors.New("query is empty")
	}

	ss := []string{}
	for key, value := range query {
		ss = append(ss, key+"="+value)
	}

	return manager.Fops(strings.Join(ss, "&"))
}
