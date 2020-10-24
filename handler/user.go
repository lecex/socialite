package handler

import (
	"context"
	"fmt"

	pb "github.com/lecex/socialite/proto/user"
	"github.com/lecex/socialite/service/repository"
)

// User 消息事件模板结构
type User struct {
	Repo repository.User
}

// Get 获取消息事件模板
func (srv *User) Get(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	socialiteUser, err := srv.Repo.Get(req.SocialiteUser)
	if err != nil {
		return err
	}
	res.SocialiteUser = socialiteUser
	return err
}

// Delete 更新消息事件模板
func (srv *User) Delete(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	valid, err := srv.Repo.Delete(req.SocialiteUser)
	if err != nil {
		res.Valid = false
		return fmt.Errorf("更新socialite配置失败")
	}
	res.Valid = valid
	return err
}
