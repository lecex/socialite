package handler

import (
	server "github.com/micro/go-micro/v2/server"

	healthPB "github.com/lecex/socialite-api/proto/health"
)

// Register 注册
func Register(Server server.Server) { //
	healthPB.RegisterHealthHandler(srv.Server, &Health{})
}
