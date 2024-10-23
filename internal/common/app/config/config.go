package config

import (
	"gopkg.in/ini.v1"
)

var App AppConfig

type Config struct {
	*ini.File
}

func New() *Config {
	cfg, err := ini.Load("./internal/common/app/config/config.ini")
	if err != nil {
		panic(err)
	}
	return &Config{
		cfg,
	}
}

type AppConfig struct {
	UsersService usersService
	PostsService postsService
	AuthService  authService
	LogService   logService
	PostgresDB   database
	MongoDB      database
}

type authService struct {
	Port string
}

type logService struct {
	Port string
}

type database struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type usersService struct {
	Port string
}

type postsService struct {
	Port string
}

func (cfg *Config) InitLog() {
	logSection := cfg.Section("log")
	App.LogService.Port = logSection.Key("port").String()
}

func (cfg *Config) InitAuth() {
	authSection := cfg.Section("auth")
	App.AuthService.Port = authSection.Key("port").String()
}

func (cfg *Config) InitUsers() {
	usersSection := cfg.Section("users")
	App.UsersService.Port = usersSection.Key("port").String()
}

func (cfg *Config) InitPosts() {
	postsSection := cfg.Section("posts")
	App.PostsService.Port = postsSection.Key("port").String()
}

func (cfg *Config) InitPostgres() {
	postgresSection := cfg.Section("postgres")
	App.PostgresDB.User = postgresSection.Key("postgres_user").String()
	App.PostgresDB.Password = postgresSection.Key("postgres_password").String()
	App.PostgresDB.Host = postgresSection.Key("postgres_host").String()
	App.PostgresDB.Port = postgresSection.Key("postgres_port").String()
	App.PostgresDB.Name = postgresSection.Key("postgres_db_name").String()
}

func (cfg *Config) InitMongo() {
	mongoSection := cfg.Section("mongo")
	App.MongoDB.User = mongoSection.Key("mongo_user").String()
	App.MongoDB.Password = mongoSection.Key("mongo_password").String()
	App.MongoDB.Host = mongoSection.Key("mongo_host").String()
	App.MongoDB.Port = mongoSection.Key("mongo_port").String()
	App.MongoDB.Name = mongoSection.Key("mongo_db_name").String()
}
