package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/lecex/socialite/handler"

	socialPB "github.com/lecex/socialite/proto/socialite"
)

func TestSocialiteAuth(t *testing.T) {
	req := &socialPB.Request{
		Socialite: &socialPB.Socialite{
			Driver: "miniprogram_wechat",
			Code:   "071mGp0w3DdE3V23KP0w37rpRM0mGp0d",
		},
	}
	res := &socialPB.Response{}
	h := handler.Socialite{}
	err := h.Auth(context.TODO(), req, res)
	fmt.Println("--------", req, res, err)
}
