package utility

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
)

// https://wcs.chinanetcenter.com/document/API/Token/AccessToken
// https://wcs.chinanetcenter.com/document/Tools/GenerateManageToken

//Auth a
type Auth struct {
	AccessKey string
	SecretKey []byte
}

//NewAuth new
func NewAuth(accessKey, secretKey string) (auth *Auth) {
	return &Auth{accessKey, []byte(secretKey)}
}

//CreateUploadToken 生成上传凭证
/// <summary>
/// https://wcs.chinanetcenter.com/document/API/Token/UploadToken
/// https://wcs.chinanetcenter.com/document/Tools/GenerateUploadToken
/// </summary>
/// <param name="putPolicy">上传策略，JSON 字符串</param>
/// <returns>上传凭证</returns>
func (a *Auth) CreateUploadToken(putPolicy string) (token string) {
	return a.SignWithData([]byte(putPolicy))
}

//Sign create sign
func (a *Auth) Sign(data []byte) (token string) {
	return a.AccessKey + ":" + a.encodeSign(data)
}

//SignWithData sign
func (a *Auth) SignWithData(data []byte) (token string) {
	encodedData := URLSafeEncode(data)
	return a.AccessKey + ":" + a.encodeSign([]byte(encodedData)) + ":" + encodedData
}

//SignRequest https://wcs.chinanetcenter.com/document/Tools/GenerateManageToken
func (a *Auth) SignRequest(reqest *http.Request) (token string, err error) {
	var data string
	u := reqest.URL
	if len(u.RawQuery) > 0 {
		data = u.Path + "?" + u.RawQuery + "\n"
	} else {
		data = u.Path + "\n"
	}

	var buffer []byte
	if reqest.Body != nil {
		var readBody io.ReadCloser
		readBody, err = reqest.GetBody()
		if nil == err {
			var body []byte
			body, err = ioutil.ReadAll(readBody)
			readBody.Close()
			if nil == err {
				buffer = append([]byte(data), body...)
			}
		}
	} else {
		buffer = []byte(data)
	}

	return a.Sign(buffer), nil
}

// private
func (a *Auth) encodeSign(data []byte) (sign string) {
	hm := hmac.New(sha1.New, a.SecretKey)
	hm.Write(data)
	sum := hm.Sum(nil)
	hexString := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(hexString, sum)
	return URLSafeEncode(hexString)
}
