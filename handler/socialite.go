package handler

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	client "github.com/lecex/core/client"
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
		fmt.Println(oauthID, err)
	case "wechat":
	default:
		err = fmt.Errorf("不支持 %s 登录", socialite.Driver)
	}

	if oauthID != "" {
		users, err = srv.getBuildUser(oauthID, socialite.Driver, req.User)
	} else {
		err = fmt.Errorf("获取授权用户 Id 失败")
	}
	fmt.Println("users**********", users, err)
	return err
}

// getBuildUser 获取绑定用户
func (srv *Socialite) getBuildUser(oauthID string, origin string, user *pb.User) (users []*pb.User, err error) {
	u := &userPB.SocialiteUser{
		OauthId: oauthID,
		Origin:  origin,
	}
	if srv.Repo.Exist(u) {
		fmt.Println(1, srv.Repo.Exist(u))
	} else {
		// 无用户先用过用户服务创建用户
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
		fmt.Println(2, req, res, err)
	}
	fmt.Println(3, u, user)

	fmt.Println("||||||||||\\\\\\\\\\\\ ", oauthID, origin, u)
	return
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
