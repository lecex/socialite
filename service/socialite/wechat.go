package socialite

import (
	"github.com/bigrocs/wechat"
	"github.com/bigrocs/wechat/requests"
	"github.com/clbanning/mxj"
)

// Wechat 微信小程序
type Wechat struct {
	AppId  string
	Secret string
}

// Oauth2 使用 code 换取 AccessToken
func (srv *Wechat) Oauth2(code string) (req mxj.Map, err error) {
	request := requests.NewCommonRequest()
	request.Domain = "offiaccount"
	request.ApiName = "sns.userinfo"
	request.QueryParams = map[string]interface{}{
		"code": code,
	}
	return srv.request(request)
}

// Userinfo 拉取用户信息(需scope为 snsapi_userinfo)
func (srv *Wechat) Userinfo(openid string, accessToken string) (req mxj.Map, err error) {
	request := requests.NewCommonRequest()
	request.Domain = "offiaccount"
	request.ApiName = "sns.oauth2.access_token"
	request.QueryParams = map[string]interface{}{
		"access_token": accessToken,
		"openid":       openid,
	}
	return srv.request(request)
}

func (srv *Wechat) request(request *requests.CommonRequest) (req mxj.Map, err error) {
	client := srv.NewClient()
	// 请求
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return req, err
	}
	req, err = response.GetHttpContentMap()
	if err != nil {
		return req, err
	}
	return req, err
}

// NewClient 创建新的连接
func (srv *Wechat) NewClient() (client *wechat.Client) {
	client = wechat.NewClient()
	c := client.Config
	c.AppId = srv.AppId
	c.Secret = srv.Secret
	return client
}
