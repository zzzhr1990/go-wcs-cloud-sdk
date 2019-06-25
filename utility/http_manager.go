package utility

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	// "github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/wcserror"
)

// 比 C# SDK 的类少了一个 allowAutoRedirect 参数，这个可以在传入的 http.Client 上自己设置 CheckRedirect
// 比 C# SDK 的类少了一个 userAgent 参数

//HTTPManager HM
type HTTPManager struct {
	//client *http.Client
}

// 这个值会在 Config 里设定
var userAgent = "WCS-GO-SDK-0.0.0.0"

//SetUserAgent set
func SetUserAgent(ua string) {
	userAgent = ua
	return
}

//CreateGetRequest cgr
func CreateGetRequest(url string) (request *CommonRequest, err error) {
	request, err = CreateCommonRequest("GET", url)
	return
}

//CreatePostRequest cpr
func CreatePostRequest(url string) (request *CommonRequest, err error) {
	request, err = CreateCommonRequest("POST", url)
	return
}

//AddMime am
func AddMime(reqest *CommonRequest, mime string) {
	if len(mime) > 0 {
		reqest.AddHeader("Content-Type", mime)
	}
}

//NewHTTPManager nhm
func NewHTTPManager() (httpManager *HTTPManager) {
	return &HTTPManager{}
}

//NewDefaultHTTPManager dnh
func NewDefaultHTTPManager() (httpManager *HTTPManager) {
	return &HTTPManager{}
}

//GetClient gpc
func (httpManager *HTTPManager) GetClient() (client *http.Client) {
	timeout := time.Duration(60 * time.Second)

	return &http.Client{
		Timeout: timeout,
		//Transport: transport,
	}

}

//GetTimeOutClient gpc
func (httpManager *HTTPManager) GetTimeOutClient(time time.Duration) (client *http.Client) {
	// timeout := time.Duration(60 * time.Second)
	if time <= 0 {
		return httpManager.GetClient()
	}
	return &http.Client{
		Timeout: time,
		//Transport: transport,
	}

}

//Do foa
func (httpManager *HTTPManager) Do(request *CommonRequest) (*http.Response, error) {
	req, err := request.CreateRequest()
	if err != nil {
		log.Errorf("cannot create request")
		return nil, err
	}
	if request.GetTimeout() > 0 {
		return httpManager.GetTimeOutClient(request.GetTimeout()).Do(req)
	}
	return httpManager.GetTimeOutClient(request.GetTimeout()).Do(req)
}

//DoRetry foa
func (httpManager *HTTPManager) DoRetry(request *CommonRequest, respEntity interface{}, retry int) error {
	for {
		req, err := request.CreateRequest()
		if err != nil {
			log.Errorf("cannot create request")
			return err
		}
		resp, err := httpManager.GetTimeOutClient(request.GetTimeout()).Do(req)
		if err == nil {
			// nil do next
			defer resp.Body.Close()
			responseBody, err := ioutil.ReadAll(resp.Body)
			// log.Infof("recv %v", string(responseBody))
			if err == nil {
				if resp.StatusCode == http.StatusOK {
					err = json.Unmarshal(responseBody, respEntity)
					if err == nil {
						return nil
					}
					log.Errorf("Http Api request Unmarshal error %v, body: %v", err, string(responseBody))
				} else {
					if resp.StatusCode == 406 {
						log.Warnf("File exists..%v", string(responseBody))
						return wcserror.ErrFileExists
					} // ErrFileNotFound
					if resp.StatusCode == 404 {
						log.Warnf("file not found..%v", req.RequestURI)
						return wcserror.ErrFileNotFound
					}
					log.Errorf("Response from API %v", string(responseBody))
					err = errors.New("Req API err")
				}
				// resp ok, json

			} else {
				log.Errorf("Request API failed,ER_FAILED_READ, %v", err)
			}

		} else {
			log.Errorf("Request API failed, %v", err)
		}
		retry--
		if retry < 1 {
			log.Info("request api give up")
			return errors.New("give up error")
		}
		//
		time.Sleep(time.Duration(2) * time.Second)
		log.Infof("retrying api %v, try time left: %v", request.uri, retry)
	}
}

//DoWithAuth wif
func (httpManager *HTTPManager) DoWithAuth(request *CommonRequest, auth *Auth) (response *http.Response, err error) {
	if nil != auth {
		request.AddAuth(auth)
	}
	return httpManager.Do(request)
}

//DoWithAuthRetry wif
func (httpManager *HTTPManager) DoWithAuthRetry(reqest *CommonRequest, auth *Auth, resp interface{}, retry int) error {
	if nil != auth {
		reqest.AddAuth(auth)
	}
	return httpManager.DoRetry(reqest, resp, retry)
}

//DoWithToken * RAW!
func (httpManager *HTTPManager) DoWithToken(request *CommonRequest, token string) (response *http.Response, err error) {
	if len(token) > 0 {
		request.AddToken(token)
	}
	return httpManager.Do(request)
}

// DoWithTokenAndRetry do Http Request With Token
func (httpManager *HTTPManager) DoWithTokenAndRetry(request *CommonRequest, token string, resp interface{}, retry int) (err error) {
	if len(token) > 0 {
		request.AddToken(token)
	}
	return httpManager.DoRetry(request, resp, retry)
}
