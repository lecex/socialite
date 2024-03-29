package repository

import (
	"fmt"
	// 公共引入

	pb "github.com/lecex/socialite/proto/user"
	"github.com/micro/go-micro/v2/util/log"

	"github.com/jinzhu/gorm"
)

//User 仓库接口
type User interface {
	GetByID(user *pb.SocialiteUser) ([]*pb.SocialiteUser, error)
	Exist(user *pb.SocialiteUser) bool
	Create(user *pb.SocialiteUser) (*pb.SocialiteUser, error)
	Get(user *pb.SocialiteUser) (*pb.SocialiteUser, error)
	Update(user *pb.SocialiteUser) (bool, error)
	UpdateByOauthId(user *pb.SocialiteUser) (bool, error)
	UpdateById(user *pb.SocialiteUser) (bool, error)
	Delete(user *pb.SocialiteUser) (bool, error)
	DeleteRelatedUser(user *pb.User) (bool, error)
}

// UserRepository 用户仓库
type UserRepository struct {
	DB *gorm.DB
}

// GetByID 根据 id 获取绑定信息
func (repo *UserRepository) GetByID(user *pb.SocialiteUser) (User []*pb.SocialiteUser, err error) {
	if user.Id != "" {
		if err := repo.DB.Model(&User).Where("id = ?", user.Id).Find(&User).Error; err != nil {
			return nil, err
		}
	}
	return User, nil
}

// Exist 检测用户是否已经存在
func (repo *UserRepository) Exist(user *pb.SocialiteUser) bool {
	var count int
	if user.Id != "" {
		repo.DB.Model(&user).Where("id = ?", user.Id).Count(&count)
		if count > 0 {
			return true
		}
	}
	if user.Origin != "" {
		repo.DB.Model(&user).Where("origin = ?", user.Origin).Where("oauth_id = ?", user.OauthId).Count(&count)
		if count > 0 {
			return true
		}
	}
	return false
}

// Get 获取用户信息
func (repo *UserRepository) Get(socialiteUser *pb.SocialiteUser) (*pb.SocialiteUser, error) {
	users := []*pb.User{}
	if socialiteUser.Id != "" {
		if err := repo.DB.Model(&socialiteUser).Where("id = ?", socialiteUser.Id).Find(&socialiteUser).Error; err != nil {
			return nil, err
		}
		if err := repo.DB.Model(&socialiteUser).Related(&users).Error; err != nil {
			if err.Error() != "record not found" {
				return nil, err
			}
		}
	}
	if socialiteUser.OauthId != "" {
		if err := repo.DB.Model(&socialiteUser).Where("origin = ?", socialiteUser.Origin).Where("oauth_id = ?", socialiteUser.OauthId).Find(&socialiteUser).Error; err != nil {
			return nil, err
		}
		if err := repo.DB.Model(&socialiteUser).Related(&users).Error; err != nil {
			if err.Error() != "record not found" {
				return nil, err
			}
		}
	}
	socialiteUser.Users = users
	return socialiteUser, nil
}

// Create 创建用户
// bug 无用户名创建用户可能引起 bug
func (repo *UserRepository) Create(user *pb.SocialiteUser) (*pb.SocialiteUser, error) {
	if exist := repo.Exist(user); exist == true {
		return user, fmt.Errorf("注册用户已存在")
	}
	err := repo.DB.Create(user).Error
	if err != nil {
		// 写入数据库未知失败记录
		fmt.Println(err)
		log.Log(err)
		return user, fmt.Errorf("注册用户失败")
	}
	return user, nil
}

// Update 更新用户
func (repo *UserRepository) Update(user *pb.SocialiteUser) (bool, error) {
	if user.Id == "" {
		return false, fmt.Errorf("未找到用户Id")
	}
	id := &pb.SocialiteUser{
		Id: user.Id,
	}
	u := &pb.SocialiteUser{
		Id:      user.Id,
		Content: user.Content,
		Users:   user.Users,
	}
	err := repo.DB.Model(id).Where("id = ?", user.Id).Updates(u).Error
	if err != nil {
		log.Log(err)
		return false, err
	}
	return true, nil
}

// UpdateByOauthId 更新用户
func (repo *UserRepository) UpdateByOauthId(user *pb.SocialiteUser) (bool, error) {
	if user.OauthId == "" {
		return false, fmt.Errorf("未找到用户OauthId")
	}
	id := &pb.SocialiteUser{
		OauthId: user.OauthId,
	}
	u := &pb.SocialiteUser{
		Id:      user.Id,
		Content: user.Content,
		Users:   user.Users,
	}
	err := repo.DB.Model(id).Where("oauth_id = ?", user.OauthId).Updates(u).Error
	if err != nil {
		log.Log(err)
		return false, err
	}
	return true, nil
}

// Update 更新用户
func (repo *UserRepository) UpdateById(user *pb.SocialiteUser) (bool, error) {
	if user.Id == "" {
		return false, fmt.Errorf("未找到用户Id")
	}
	id := &pb.SocialiteUser{
		Id: user.Id,
	}
	u := &pb.SocialiteUser{
		Id:      user.Id,
		OauthId: user.OauthId,
		Content: user.Content,
		Users:   user.Users,
	}
	err := repo.DB.Model(id).Where("id = ?", user.Id).Updates(u).Error
	if err != nil {
		log.Log(err)
		return false, err
	}
	return true, nil
}

// Delete 删除用户
func (repo *UserRepository) Delete(user *pb.SocialiteUser) (bool, error) {
	if user.Id == "" {
		return false, fmt.Errorf("请传入更新id")
	}
	err := repo.DB.Where("id = ?", user.Id).Delete(user).Error
	if err != nil {
		log.Log(err)
		return false, err
	}
	return true, nil
}

// DeleteRelatedUser 删除关联用户信息
func (repo *UserRepository) DeleteRelatedUser(user *pb.User) (bool, error) {
	if user.Id == "" {
		return false, fmt.Errorf("请传入更新id")
	}
	err := repo.DB.Where("id = ?", user.Id).Delete(user).Error
	if err != nil {
		log.Log(err)
		return false, err
	}
	return true, nil
}
