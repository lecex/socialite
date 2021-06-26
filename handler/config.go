package handler

import (
	"context"
	"fmt"

	pb "github.com/lecex/socialite/proto/config"
	"github.com/lecex/socialite/service/repository"
)

// Config 消息事件模板结构
type Config struct {
	Repo repository.Config
}

// All 获取所有权限
func (srv *Config) All(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	configs, err := srv.Repo.All(req)
	if err != nil {
		return err
	}
	res.Configs = configs
	return err
}

// List 获取所有权限
func (srv *Config) List(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	configs, err := srv.Repo.List(req.ListQuery)
	total, err := srv.Repo.Total(req.ListQuery)
	if err != nil {
		return err
	}
	res.Configs = configs
	res.Total = total
	return err
}

// Get 获取权限
func (srv *Config) Get(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	config, err := srv.Repo.Get(req.Config)
	if err != nil {
		return err
	}
	res.Config = config
	return err
}

// Create 创建权限
func (srv *Config) Create(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	_, err = srv.Repo.Create(req.Config)
	if err != nil {
		res.Valid = false
		return fmt.Errorf("添加权限失败")
	}
	res.Valid = true
	return err
}

// Update 更新权限
func (srv *Config) Update(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	valid, err := srv.Repo.Update(req.Config)
	if err != nil {
		res.Valid = false
		return fmt.Errorf("更新权限失败")
	}
	res.Valid = valid
	return err
}

// Delete 删除权限
func (srv *Config) Delete(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	valid, err := srv.Repo.Delete(req.Config)
	if err != nil {
		res.Valid = false
		return fmt.Errorf("删除权限失败")
	}
	res.Valid = valid
	return err
}
