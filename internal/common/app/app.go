package app

import (
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"testing"
)

func InitUsersService() {
	config.InitUsersConfig()
	log.Init()
}

func InitMock(c config.AppConfig, t testing.TB){
	config.App = c
	log.InitMock(t)
}