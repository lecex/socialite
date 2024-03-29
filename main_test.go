package main

import (
	"context"
	"fmt"
	"testing"

	db "github.com/lecex/socialite/providers/database"

	"github.com/lecex/socialite/config"
	"github.com/lecex/socialite/handler"

	conPB "github.com/lecex/socialite/proto/config"

	_ "github.com/lecex/socialite/providers/migrations" // 执行数据迁移
	"github.com/lecex/socialite/service/repository"

	socialPB "github.com/lecex/socialite/proto/socialite"
)

var Conf = config.Conf

func TestSocialiteConfigUpdate(t *testing.T) {
	req := &conPB.Request{
		Config: &conPB.Config{
			Name:         "微信小程序",
			Driver:       "miniprogram_wechat",
			ClientId:     "wx15550c1a89d982c8",
			ClientSecret: "f9c11f183a5beb592ccd801298ff5533",
			Status:       true,
		},
	}
	res := &conPB.Response{}
	h := handler.Config{&repository.ConfigRepository{db.DB}}
	err := h.Create(context.TODO(), req, res)
	fmt.Println("--------")
	t.Log(req, res, err)
}

func TestSocialiteAuth(t *testing.T) {
	req := &socialPB.Request{
		Socialite: &socialPB.Socialite{
			Driver:   "miniprogram_wechat",
			ClientId: "wx15550c1a89d982c8",
			Code:     "051keH000YIkWL1ld5300T6JAP1keH0R",
		},
	}
	res := &socialPB.Response{}
	h := handler.Socialite{
		&repository.UserRepository{db.DB},
		&repository.ConfigRepository{db.DB},
		Conf.Service["user"],
	}
	err := h.Auth(context.TODO(), req, res)
	fmt.Println("----Auth----", res, err)
}

// func TestSocialiteBuildUser(t *testing.T) {
// 	req := &socialPB.Request{
// 		SocialiteUser: &socialPB.SocialiteUser{
// 			Id: "e39d3d2f-a978-4a70-8683-b53379e52ec1",
// 			Users: []*socialPB.User{
// 				{
// 					Id: "c9d7a398-4d59-480e-b435-469365307f8c",
// 				},
// 			},
// 		},
// 	}
// 	res := &socialPB.Response{}
// 	h := handler.Socialite{
// 		&repository.UserRepository{db.DB},
// 		Conf.Service["user"],
// 	}
// 	err := h.BuildUser(context.TODO(), req, res)
// 	fmt.Println("---Register---", req, "||", res, err)
// }
