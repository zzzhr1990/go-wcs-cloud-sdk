package utility

import (
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"

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
func CreateGetRequest(url string) (request *http.Request, err error) {
	request, err = http.NewRequest("GET", url, nil)
	return
}

//CreatePostRequest cpr
func CreatePostRequest(url string) (request *http.Request, err error) {
	request, err = http.NewRequest("POST", url, nil)
	return
}

//AddMime am
func AddMime(reqest *http.Request, mime string) {
	if len(mime) > 0 {
		reqest.Header.Set("Content-Type", mime)
	}
	return
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
	timeout := time.Duration(30 * time.Second)

	return &http.Client{
		Timeout: timeout,
		//Transport: transport,
	}

}

//Do foa
func (httpManager *HTTPManager) Do(reqest *http.Request) (response *http.Response, err error) {
	if _, ok := reqest.Header["User-Agent"]; !ok {
		reqest.Header.Set("User-Agent", userAgent)
	}
	return httpManager.GetClient().Do(reqest)
}

//DoRetry foa
func (httpManager *HTTPManager) DoRetry(reqest *http.Request, resp interface{}, retry int) (err error) {
	if _, ok := reqest.Header["User-Agent"]; !ok {
		reqest.Header.Set("User-Agent", userAgent)
	}
	// return httpManager.GetClient().Do(reqest)
	for {
		resp, err := httpManager.GetClient().Do(reqest)
		if err == nil{
			// nil do next 
			defer resp.Body.Close()
			responseBody, err := ioutil.ReadAll(resp.Body)
			if err == nil{
				if resp.StatusCode == http.StatusOK {
					log.Infof("Http Api request OKKKKKKKK %v", string(responseBody))
					err = json.Unmarshal(responseBody, resp)
					if err == nil{
						return nil
					} 
					log.Errorf("Http Api request Unmarshal error %v", err)
				}else{
					log.Errorf("Response from API %v", string(responseBody))
					err = errors.New("Req API err")
				}
				// resp ok, json
				
			}else{
				log.Errorf("Request API failed,ER_FAILED_READ, %v", err)
			}

		}else{
			log.Errorf("Request API failed, %v", err)
		}
		retry--
		if retry < 1{
			return err
		}
	}
}

//DoWithAuth wif
func (httpManager *HTTPManager) DoWithAuth(reqest *http.Request, auth *Auth) (response *http.Response, err error) {
	if nil != auth {
		var token string
		token, err = auth.SignRequest(reqest)
		if nil != err {
			return
		}
		reqest.Header.Set("Authorization", token)
	}
	return httpManager.Do(reqest)
}

//DoWithAuthRetry wif
func (httpManager *HTTPManager) DoWithAuthRetry(reqest *http.Request, auth *Auth, resp interface{}, retry int) (error) {
	if nil != auth {
		token, err := auth.SignRequest(reqest)
		if nil != err {
			return err
		}
		reqest.Header.Set("Authorization", token)
	}
	return httpManager.DoRetry(reqest, resp, retry)
}

//DoWithToken dej
func (httpManager *HTTPManager) DoWithToken(reqest *http.Request, token string) (response *http.Response, err error) {
	if len(token) > 0 {
		reqest.Header.Set("Authorization", token)
	}
	return httpManager.Do(reqest)
}
