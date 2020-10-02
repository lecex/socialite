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
			Code:   "0011Wc000kc7oK1VQx100mU2bM11Wc0G",
		},
	}
	res := &socialPB.Response{}
	h := handler.Socialite{}
	err := h.Auth(context.TODO(), req, res)
	fmt.Println("--------", req, res, err)
}
