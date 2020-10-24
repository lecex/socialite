package handler

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/clbanning/mxj"
	// client "github.com/lecex/core/client"
	// authSrvPB "github.com/lecex/user/proto/auth"
	// userSrvPB "github.com/lecex/user/proto/user"

	conPB "github.com/lecex/socialite/proto/config"
	pb "github.com/lecex/socialite/proto/socialite"
	userPB "github.com/lecex/socialite/proto/user"

	db "github.com/lecex/socialite/providers/database"
	"github.com/lecex/socialite/service/repository"
	"github.com/lecex/socialite/service/socialite"
)

// Socialite 社会登录
type Socialite struct {
	Repo        repository.User
	ServiceName string
}

// Auth 登录授权
func (srv *Socialite) Auth(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	var content mxj.Map
	socialite := req.Socialite
	switch socialite.Driver {
	case "miniprogram_wechat":
		content, err = srv.miniprogramWechat(socialite.Code)
	case "wechat":
		content, err = srv.wechat(socialite.Code)
	default:
		err = fmt.Errorf("不支持 %s 登录", socialite.Driver)
	}
	if err != nil {
		return err
	}
	if _, ok := content["oauthid"]; ok {
		// 获取相关用户信息
		res.SocialiteUser, err = srv.getSocialiteUser(content, socialite.Driver)
	} else {
		err = fmt.Errorf("获取授权用户 Id 失败")
	}
	return err
}

// Register 注册
// func (srv *Socialite) Register(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
// 	u := &userPB.SocialiteUser{
// 		Id: req.SocialiteUser.Id,
// 	}
// 	// 获取所有关联用户
// 	u, err = srv.Repo.Get(u)
// 	if err != nil {
// 		return err
// 	}
// 	if len(req.SocialiteUser.Users) > 0 {
// 		for _, user := range req.SocialiteUser.Users {
// 			// 无用户先通过用户服务创建用户
// 			reqUserSrv := &userSrvPB.Request{
// 				User: &userSrvPB.User{
// 					Username: user.Username,
// 					Mobile:   user.Mobile,
// 					Email:    user.Email,
// 					Password: user.Password,
// 					Name:     user.Name,
// 					Avatar:   user.Avatar,
// 				},
// 			}
// 			resUserSrv := &userSrvPB.Response{}
// 			err = client.Call(context.TODO(), srv.ServiceName, "Users.Create", reqUserSrv, resUserSrv)
// 			if err != nil {
// 				return err
// 			}
// 			if resUserSrv.Valid {
// 				u.Users = append(u.Users, &userPB.User{
// 					Id: resUserSrv.User.Id,
// 				})
// 			}
// 		}
// 	} else {
// 		err = fmt.Errorf("未收到用户注册信息")
// 	}
// 	u.CreatedAt = ""
// 	u.UpdatedAt = ""
// 	_, err = srv.Repo.Update(u)
// 	fmt.Println("---Register---", u)
// 	return err
// }

// AuthURL 授权网址
func (srv *Socialite) AuthURL(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return err
}

// getSocialiteUser 获取绑定用户
func (srv *Socialite) getSocialiteUser(content mxj.Map, origin string) (socialiteUser *pb.SocialiteUser, err error) {
	c, _ := content.Json()
	u := &userPB.SocialiteUser{
		OauthId: content["oauthid"].(string),
		Origin:  origin,
		Content: string(c),
	}
	if srv.Repo.Exist(u) {
		_, err = srv.Repo.Update(u)
		if err != nil {
			return nil, err
		}
	} else {
		_, err = srv.Repo.Create(u)
		if err != nil {
			return nil, err
		}
	}
	// 获取所有关联用户
	u, err = srv.Repo.Get(u)
	if err != nil {
		return nil, err
	}
	socialiteUser = &pb.SocialiteUser{
		Id:      u.Id,
		OauthId: u.OauthId,
		Origin:  u.Origin,
		Content: u.Content,
	}
	// 获取关联用户token
	for _, user := range u.Users {
		// reqAuthSrv := &authSrvPB.Request{
		// 	User: &authSrvPB.User{
		// 		Id: user.Id,
		// 	},
		// }
		// resAuthSrv := &authSrvPB.Response{}
		// err = client.Call(context.TODO(), srv.ServiceName, "Auth.AuthById", reqAuthSrv, resAuthSrv)
		// if err != nil {
		// 	return nil, err
		// }
		socialiteUser.Users = append(socialiteUser.Users, &pb.User{
			Id: user.Id,
			// Name:  resAuthSrv.User.Name,
			// Token: resAuthSrv.Token,
		})
	}
	return socialiteUser, nil
}

// getConfig 初始化配置等
func (srv *Socialite) getConfig() (*conPB.Config, error) {
	res := &conPB.Response{}
	h := Config{&repository.ConfigRepository{db.DB}}
	err := h.Get(context.TODO(), &conPB.Request{}, res)
	if res.Config == nil {
		err = fmt.Errorf("获取配置失败")
	}
	return res.Config, err
}

// 微信小程序 获取 oauthID
func (srv *Socialite) miniprogramWechat(code string) (req mxj.Map, err error) {
	con, err := srv.getConfig()
	if err != nil {
		return
	}
	if con.MiniprogramWechat == nil {
		err = fmt.Errorf("未配置微信小程序")
		return
	}
	m := &socialite.MiniprogramWechat{
		AppId:  con.MiniprogramWechat.AppId,
		Secret: con.MiniprogramWechat.Secret,
	}
	req, err = m.Code2Session(code)
	if err != nil {
		return
	}
	if _, ok := req["errmsg"]; ok {
		err = fmt.Errorf(req["errmsg"].(string))
	}
	if _, ok := req["openid"]; ok {
		req["oauthid"] = req["openid"]
	}
	if _, ok := req["unionid"]; ok {
		req["oauthid"] = req["unionid"]
	}
	return
}

// wechat 微信获取 oauthID
func (srv *Socialite) wechat(code string) (req mxj.Map, err error) {
	con, err := srv.getConfig()
	if err != nil {
		return
	}
	m := &socialite.Wechat{
		AppId:  con.Wechat.AppId,
		Secret: con.Wechat.Secret,
	}
	req, err = m.Oauth2(code)
	if err != nil {
		return
	}
	if _, ok := req["errmsg"]; ok {
		err = fmt.Errorf(req["errmsg"].(string))
	}
	if _, ok := req["openid"]; ok {
		req["oauthid"] = req["openid"]
	}
	if _, ok := req["unionid"]; ok {
		req["oauthid"] = req["unionid"]
	}
	return
}

// getRandomString 生成随机字符串
func (srv *Socialite) getRandomString(length int64) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; int64(i) < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
