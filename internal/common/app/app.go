package app

import (
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"testing"
)

func defaultInit(cfg *config.Config) {
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

func InitProjectsService() {
	cfg := config.New()
	defaultInit(cfg)
	cfg.InitProjects()
	cfg.InitPostgres()
	cfg.InitUsers()
}

func InitPostsService() {
	cfg := config.New()
	defaultInit(cfg)
	cfg.InitPosts()
	cfg.InitFiles()
	cfg.InitUsers()
	cfg.InitChats()
	cfg.InitMongo()
}

func InitChatsService() {
	cfg := config.New()
	defaultInit(cfg)
	cfg.InitChats()
	cfg.InitUsers()
	cfg.InitPostgres()
}

func InitAuthService() {
	cfg := config.New()
	defaultInit(cfg)
}

func InitLogService() {
	cfg := config.New()
	cfg.InitLog()
}

func InitFilesService() {
	cfg := config.New()
	defaultInit(cfg)
	cfg.InitFiles()
}

func InitMock(c config.AppConfig, t testing.TB) {
	config.App = c
	log.InitMock(t)
}
