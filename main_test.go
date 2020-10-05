package main

import (
	"fmt"
	"testing"

	"github.com/lecex/socialite/config"
	// 	"context"
	// 	"fmt"
	// 	"testing"
	// 	db "github.com/lecex/socialite/providers/database"
	// 	// _ "github.com/lecex/socialite/providers/migrations" // 执行数据迁移
	// 	"github.com/lecex/socialite/config"
	// 	"github.com/lecex/socialite/handler"
	// 	conPB "github.com/lecex/socialite/proto/config"
	// 	socialPB "github.com/lecex/socialite/proto/socialite"
	// 	"github.com/lecex/socialite/service/repository"
)

var Conf = config.Conf

func TestSocialiteConfigUpdate(t *testing.T) {
	// req := &conPB.Request{
	// 	Config: &conPB.Config{
	// 		MiniprogramWechat: &conPB.MiniprogramWechat{
	// 			AppId:  "wx15550c1a89d982c8",
	// 			Secret: "f9c11f183a5beb592ccd801298ff5533",
	// 		},
	// 	},
	// }
	// res := &conPB.Response{}
	// h := handler.Config{&repository.ConfigRepository{db.DB}}
	// err := h.Update(context.TODO(), req, res)
	fmt.Println("--------")
	// t.Log(req, res, err)
}

// func TestSocialiteAuth(t *testing.T) {

// 	fmt.Println(Conf)
// 	req := &socialPB.Request{
// 		Socialite: &socialPB.Socialite{
// 			Driver: "miniprogram_wechat",
// 			Code:   "071ZNlGa1ETYJz0GvPFa1x4UaC4ZNlGh",
// 		},
// 		User: &socialPB.User{
// 			Name:   "BigRocs",
// 			Avatar: "https://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83ep1m5aI7y3WJAP6XIXN4e39124xvcjJoI9AM8QXjB9jN6VJpl3C32VNeXELnB71EWk8sE7zp32n4A/132",
// 		},
// 	}
// 	res := &socialPB.Response{}
// 	h := handler.Socialite{
// 		&repository.UserRepository{db.DB},
// 		Conf.Service["user"],
// 	}
// 	err := h.Auth(context.TODO(), req, res)
// 	fmt.Println("--------", req, res, err)
// }
