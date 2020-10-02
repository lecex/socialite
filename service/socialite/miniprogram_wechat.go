package socialite

import (
	"fmt"

	"github.com/bigrocs/wechat"
	"github.com/bigrocs/wechat/requests"
	"github.com/clbanning/mxj"
)

// MiniprogramWechat 微信小程序
type MiniprogramWechat struct {
	AppId      string
	Secret     string
	SessionKey string
}

// Code2Session 使用 code 换取 session
func (srv *MiniprogramWechat) Code2Session(js_code string) (req mxj.Map, err error) {
	request := requests.NewCommonRequest()
	request.Domain = "miniprogram"
	request.ApiName = "auth.code2Session"
	request.QueryParams = map[string]interface{}{
		"js_code": js_code,
	}
	return srv.request(request)
}

func (srv *MiniprogramWechat) request(request *requests.CommonRequest) (req mxj.Map, err error) {
	client := srv.NewClient()
	// 请求
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return req, err
	}
	req, err = response.GetHttpContentMap()
	fmt.Println(client, req)
	if err != nil {
		return req, err
	}
	return req, err
}

// NewClient 创建新的连接
func (srv *MiniprogramWechat) NewClient() (client *wechat.Client) {
	client = wechat.NewClient()
	c := client.Config
	c.AppId = srv.AppId
	c.Secret = srv.Secret
	// c.SessionKey = srv.SessionKey
	return client
}
