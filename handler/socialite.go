package handler

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	client "github.com/lecex/core/client"
	authSrvPB "github.com/lecex/user/proto/auth"
	userSrvPB "github.com/lecex/user/proto/user"

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

// Auth 小程序登录授权
func (srv *Socialite) Auth(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	var users []*pb.User
	var oauthID string
	socialite := req.Socialite
	switch socialite.Driver {
	case "miniprogram_wechat":
		oauthID, err = srv.miniprogramWechat(socialite.Code)
	case "wechat":
	default:
		err = fmt.Errorf("不支持 %s 登录", socialite.Driver)
	}
	if oauthID != "" {
		users, err = srv.getBuildUser(oauthID, socialite.Driver, req.User)
	} else {
		err = fmt.Errorf("获取授权用户 Id 失败")
	}
	res.Users = users
	return err
}

// getBuildUser 获取绑定用户
func (srv *Socialite) getBuildUser(oauthID string, origin string, user *pb.User) (users []*pb.User, err error) {
	u := &userPB.SocialiteUser{
		OauthId: oauthID,
		Origin:  origin,
	}
	if !srv.Repo.Exist(u) {
		// 无用户先通过用户服务创建用户
		req := &userSrvPB.Request{
			User: &userSrvPB.User{
				Name:     user.Name,
				Avatar:   user.Avatar,
				Origin:   u.Origin,
				Password: srv.getRandomString(16), // 密码默认为 16 位随机数
			},
		}
		res := &userSrvPB.Response{}
		err = client.Call(context.TODO(), srv.ServiceName, "Users.Create", req, res)
		if err != nil {
			return nil, err
		}
		u.Users = append(u.Users, &userPB.User{
			Id:   res.User.Id,
			Name: res.User.Name,
		})
		_, err = srv.Repo.Create(u)
		if err != nil {
			return nil, err
		}
	}
	// 获取所有关联用户
	socialiteUser, err := srv.Repo.Get(u)
	if err != nil {
		return nil, err
	}
	// 获取关联用户token
	for _, user := range socialiteUser.Users {
		// 无用户先通过用户服务创建用户
		req := &authSrvPB.Request{
			User: &authSrvPB.User{
				Id: user.Id,
			},
		}
		res := &authSrvPB.Response{}
		err = client.Call(context.TODO(), srv.ServiceName, "Auth.AuthById", req, res)
		if err != nil {
			return nil, err
		}
		users = append(users, &pb.User{
			Id:    user.Id,
			Name:  user.Name,
			Token: res.Token,
		})
	}
	return users, nil
}

// getConfig 初始化配置等
func (srv *Socialite) getConfig() (*conPB.Config, error) {
	res := &conPB.Response{}
	h := Config{&repository.ConfigRepository{db.DB}}
	err := h.Get(context.TODO(), &conPB.Request{}, res)
	return res.Config, err
}

func (srv *Socialite) miniprogramWechat(code string) (oauthID string, err error) {
	con, err := srv.getConfig()
	if err != nil {
		return
	}
	m := &socialite.MiniprogramWechat{
		AppId:  con.MiniprogramWechat.AppId,
		Secret: con.MiniprogramWechat.Secret,
	}
	req, err := m.Code2Session(code)
	if err != nil {
		return
	}
	if _, ok := req["errmsg"]; ok {
		err = fmt.Errorf(req["errmsg"].(string))
	}
	if _, ok := req["openid"]; ok {
		oauthID = req["openid"].(string)
	}
	if _, ok := req["unionid"]; ok {
		oauthID = req["unionid"].(string)
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
