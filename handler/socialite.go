package handler

import (
	"context"
	"fmt"

	pb "github.com/lecex/socialite/proto/socialite"
	"github.com/lecex/socialite/service/socialite"
)

// Socialite 社会登录
type Socialite struct {
}

// Auth 小程序登录授权
func (srv *Socialite) Auth(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	var users []*pb.User
	socialite := req.Socialite
	switch socialite.Driver {
	case "miniprogram_wechat":
		users, err = srv.miniprogramWechat(socialite.Code)
	case "wechat":
	default:
		err = fmt.Errorf("不支持 %s 登录", socialite.Driver)
	}
	fmt.Println(users, err)
	return err
}

func (srv *Socialite) miniprogramWechat(code string) (users []*pb.User, err error) {
	m := &socialite.MiniprogramWechat{
		AppId:      "wx15550c1a89d982c8",
		Secret:     "f9c11f183a5beb592ccd801298ff5533",
		SessionKey: "",
	}
	req, err := m.Code2Session(code)
	if err != nil {
		return nil, err
	}
	fmt.Println(code, req, err)
	return
}
