package repository

import (

	// 公共引入

	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lecex/core/util"
	pb "github.com/lecex/socialite/proto/config"
	"github.com/micro/go-micro/v2/util/log"
)

//Config 仓库接口
type Config interface {
	Create(config *pb.Config) (*pb.Config, error)
	Delete(config *pb.Config) (bool, error)
	Update(config *pb.Config) (bool, error)
	Get(config *pb.Config) (*pb.Config, error)
	GetByClientId(config *pb.Config) (*pb.Config, error)
	All(req *pb.Request) ([]*pb.Config, error)
	List(req *pb.ListQuery) ([]*pb.Config, error)
	Total(req *pb.ListQuery) (int64, error)
}

// ConfigRepository 模版仓库
type ConfigRepository struct {
	DB *gorm.DB
}

// All 获取所有角色信息
func (repo *ConfigRepository) All(req *pb.Request) (configs []*pb.Config, err error) {
	if err := repo.DB.Find(&configs).Error; err != nil {
		log.Log(err)
		return nil, err
	}
	return configs, nil
}

// List 获取所有权限信息
func (repo *ConfigRepository) List(req *pb.ListQuery) (configs []*pb.Config, err error) {
	db := repo.DB
	limit, offset := util.Page(req.Limit, req.Page) // 分页
	sort := util.Sort(req.Sort)                     // 排序 默认 created_at desc
	// 查询条件
	if req.Where != "" {
		db = db.Where(req.Where)
	}
	if err := db.Order(sort).Limit(limit).Offset(offset).Find(&configs).Error; err != nil {
		log.Log(err)
		return nil, err
	}
	return configs, nil
}

// Total 获取所有权限查询总量
func (repo *ConfigRepository) Total(req *pb.ListQuery) (total int64, err error) {
	configs := []pb.Config{}
	db := repo.DB
	// 查询条件
	if req.Where != "" {
		db = db.Where(req.Where)
	}
	if err := db.Find(&configs).Count(&total).Error; err != nil {
		log.Log(err)
		return total, err
	}
	return total, nil
}

// Get 获取权限信息
func (repo *ConfigRepository) Get(config *pb.Config) (*pb.Config, error) {
	if err := repo.DB.Where(&config).Find(&config).Error; err != nil {
		return nil, err
	}
	return config, nil
}

// Get 获取权限信息
func (repo *ConfigRepository) GetByClientId(config *pb.Config) (*pb.Config, error) {
	if err := repo.DB.Where("driver = ?", config.Driver).Where("client_id = ?", config.ClientId).Find(&config).Error; err != nil {
		return nil, err
	}
	return config, nil
}

// Create 创建权限
// bug 无权限名创建权限可能引起 bug
func (repo *ConfigRepository) Create(p *pb.Config) (*pb.Config, error) {
	err := repo.DB.Create(p).Error
	if err != nil {
		// 写入数据库未知失败记录
		log.Log(err)
		return p, fmt.Errorf("添加权限失败")
	}
	return p, nil
}

// Update 更新权限
func (repo *ConfigRepository) Update(p *pb.Config) (bool, error) {
	if p.Id == 0 {
		return false, fmt.Errorf("请传入更新id")
	}
	id := &pb.Config{
		Id: p.Id,
	}
	err := repo.DB.Model(id).Where("id = ?", id.Id).Updates(p).Error
	if err != nil {
		log.Log(err)
		return false, err
	}
	return true, nil
}

// Delete 删除权限
func (repo *ConfigRepository) Delete(p *pb.Config) (bool, error) {
	err := repo.DB.Delete(p).Error
	if err != nil {
		log.Log(err)
		return false, err
	}
	return true, nil
}
