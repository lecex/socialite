package handler

import (
	"context"

	pb "github.com/lecex/socialite/proto/socialite"
)

// Socialite 社会登录
type Socialite struct {
}

// Auth 小程序登录授权
func (srv *Socialite) Auth(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return err
}
