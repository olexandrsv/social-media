package app

import (
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"testing"
)

func defaultInit(cfg *config.Config){
	cfg.InitAuth()
	cfg.InitLog()
	log.Init()
}

func InitUsersService() {
	cfg := config.New()
	defaultInit(cfg)
	cfg.InitUsers()
	cfg.InitPostgres()
}

func InitPostsService(){
	cfg := config.New()
	defaultInit(cfg)
	cfg.InitPosts()
	cfg.InitMongo()
}

func InitMock(c config.AppConfig, t testing.TB){
	config.App = c
	log.InitMock(t)
}