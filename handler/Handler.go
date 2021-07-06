package handler

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"

	configPB "github.com/lecex/socialite/proto/config"
	socialitePB "github.com/lecex/socialite/proto/socialite"
	userPB "github.com/lecex/socialite/proto/user"

	db "github.com/lecex/socialite/providers/database"
	service "github.com/lecex/socialite/service/repository"
)

const topic = "event"

// Register 注册
func Register(srv micro.Service) {
	server := srv.Server()
	configPB.RegisterConfigsHandler(server, &Config{&service.ConfigRepository{db.DB}})
	socialitePB.RegisterSocialitesHandler(server, &Socialite{
		&service.UserRepository{db.DB},
		&service.ConfigRepository{db.DB},
		"user",
	})
	userPB.RegisterUsersHandler(server, &User{
		&service.UserRepository{db.DB},
	})
	// 事件处理
	sub := &Subscriber{&service.UserRepository{db.DB}}
	err := micro.RegisterSubscriber(topic, server, sub)
	if err != nil {
		log.Fatal(err)
	}
}
