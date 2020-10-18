package handler

import (
	"context"
	"fmt"

	pb "github.com/lecex/socialite/proto/User"
	"github.com/lecex/socialite/service/repository"
)

// User 消息事件模板结构
type User struct {
	Repo repository.User
}

// Get 获取消息事件模板
func (srv *User) Get(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	nat, err := srv.Repo.Get()
	if err != nil {
		return err
	}
	res.User = nat
	return err
}

// Update 更新消息事件模板
func (srv *User) Update(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	valid, err := srv.Repo.Update(req.User)
	if err != nil {
		res.Valid = false
		return fmt.Errorf("更新socialite配置失败")
	}
	res.Valid = valid
	return err
}
