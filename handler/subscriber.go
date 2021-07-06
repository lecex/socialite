package handler

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/v2/util/log"

	pb "github.com/lecex/core/proto/event"
	userPB "github.com/lecex/socialite/proto/user"

	"github.com/lecex/socialite/service/repository"
)

// Subscriber 结构
type Subscriber struct {
	Repo repository.User
}

// Process 事件处理
func (sub *Subscriber) Process(ctx context.Context, event *pb.Event) (err error) {
	switch event.Topic {
	case "user.Users.Delete":
		sub.delete(ctx, event.Data)
	}
	return err
}

// delete 删除
func (sub *Subscriber) delete(ctx context.Context, data []byte) (err error) {
	r := make(map[string]interface{})
	err = json.Unmarshal(data, &r)
	if err != nil {
		return
	}
	// 获取设备信息
	if userID, ok := r["id"]; ok {
		user := &userPB.User{
			Id: userID.(string),
		}
		valid, err := sub.Repo.DeleteRelatedUser(user)
		if err != nil {
			log.Fatal("删除社会登录账号失败", valid, err)
		}
	}
	return
}
