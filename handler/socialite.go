package handler

import (
	"context"
	"fmt"

	"github.com/clbanning/mxj"
	// client "github.com/lecex/core/client"
	// authSrvPB "github.com/lecex/user/proto/auth"
	// userSrvPB "github.com/lecex/user/proto/user"

	conPB "github.com/lecex/socialite/proto/config"
	pb "github.com/lecex/socialite/proto/socialite"
	userPB "github.com/lecex/socialite/proto/user"

	"github.com/lecex/socialite/service/repository"
	"github.com/lecex/socialite/service/socialite"
)

// Socialite 社会登录
type Socialite struct {
	Repo        repository.User
	ConfRepo    repository.Config
	ServiceName string
}

// Auth 登录授权
func (srv *Socialite) Auth(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	var content mxj.Map
	socialite := req.Socialite
	con, err := srv.getConfig(socialite)
	if err != nil {
		return err
	}
	switch socialite.Driver {
	case "miniprogram_wechat":
		content, err = srv.miniprogramWechat(socialite.Code, con)
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

// AuthURL 授权网址
func (srv *Socialite) AuthURL(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return err
}

// getConfig 初始化配置等
func (srv *Socialite) getConfig(socialite *pb.Socialite) (*conPB.Config, error) {
	con := &conPB.Config{
		Driver:   socialite.Driver,
		ClientId: socialite.ClientId,
	}
	return srv.ConfRepo.GetByClientId(con)
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
		_, err = srv.Repo.UpdateByOauthId(u)
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
		socialiteUser.Users = append(socialiteUser.Users, &pb.User{
			Id: user.Id,
		})
	}
	return socialiteUser, nil
}

// 微信小程序 获取 oauthID
func (srv *Socialite) miniprogramWechat(code string, con *conPB.Config) (req mxj.Map, err error) {
	m := &socialite.MiniprogramWechat{
		AppId:  con.ClientId,
		Secret: con.ClientSecret,
	}
	req, err = m.Code2Session(code)
	if err != nil {
		return
	}
	if _, ok := req["errmsg"]; ok {
		err = fmt.Errorf(req["errmsg"].(string))
	}
	// 默认openid作为唯一识别
	if _, ok := req["openid"]; ok {
		req["oauthid"] = req["openid"]
	}
	// 如果unionid存在在更新 openid 为 unionid 作为唯一识别
	if _, ok := req["unionid"]; ok {
		u := &userPB.SocialiteUser{
			OauthId: req["openid"].(string),
			Origin:  con.Driver,
		}
		if srv.Repo.Exist(u) {
			u, err = srv.Repo.Get(u)
			if err != nil {
				return nil, err
			}
			u.OauthId = req["unionid"].(string)
			_, err = srv.Repo.Update(u)
			if err != nil {
				return nil, err
			}
		}
		req["oauthid"] = req["unionid"]
	}
	return
}

// BuildUser 绑定用户
// Auth 登录授权
func (srv *Socialite) BuildUser(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	u := &userPB.SocialiteUser{
		Id: req.SocialiteUser.Id,
		Users: []*userPB.User{
			{
				Id:              req.SocialiteUser.Users[0].Id,
				SocialiteUserId: req.SocialiteUser.Id,
			},
		},
	}
	_, err = srv.Repo.Update(u)
	if err != nil {
		res.Valid = false
		return err
	}
	res.Valid = true
	return
}
