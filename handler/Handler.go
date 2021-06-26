package handler

import (
	server "github.com/micro/go-micro/v2/server"

	configPB "github.com/lecex/socialite/proto/config"
	socialitePB "github.com/lecex/socialite/proto/socialite"
	userPB "github.com/lecex/socialite/proto/user"

	db "github.com/lecex/socialite/providers/database"
	"github.com/lecex/socialite/service/repository"
)

// Register 注册
func Register(Server server.Server) {
	configPB.RegisterConfigsHandler(Server, &Config{&repository.ConfigRepository{db.DB}})
	socialitePB.RegisterSocialitesHandler(Server, &Socialite{
		&repository.UserRepository{db.DB},
		&repository.ConfigRepository{db.DB},
		"user",
	})
	userPB.RegisterUsersHandler(Server, &User{
		&repository.UserRepository{db.DB},
	})
}
