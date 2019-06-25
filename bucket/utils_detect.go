package bucket

import (
	"errors"

	"strings"

	"github.com/zzzhr1990/go-wcs-cloud-sdk/core"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
)

// ImageDetect Detect detect Porns? https://wcs.chinanetcenter.com/document/API/Image-op/imageDetect
func (manager *Manager) ImageDetect(imageURL string, bucket string, detectType string) (*core.DetectResponse, error) {
	if 0 == len(imageURL) {
		return nil, errors.New("query is empty")
	}

	url := manager.config.GetManageURLPrefix() + "/imageDetect"
	request, err := utility.CreatePostRequest(url)
	// log.Printf("POST %v", url)
	if nil != err {
		return nil, err
	}

	var sb strings.Builder
	sb.WriteString("type=")
	sb.WriteString(detectType)
	sb.WriteString("&image=")
	sb.WriteString(utility.URLSafeEncodeString(imageURL))
	sb.WriteString("&bucket=")
	sb.WriteString(bucket)
	// type=porn&image=aHR0cDovL3d3dy5iYWlkdS5jb20=&bucket=bucketName

	request.AddStringBody(sb.String())
	// log.Infof("call fops %v", query)
	respEntity := &core.DetectResponse{}
	err = manager.httpManager.DoWithAuthRetry(request, manager.auth, respEntity, 10)
	if err != nil {
		return nil, err
	}
	return respEntity, nil
}
