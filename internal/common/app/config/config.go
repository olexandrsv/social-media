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
	Users        usersService
	PostsService postsService
	Chats        chatsService
	Projects     projectsService
	AuthService  authService
	LogService   logService
	FilesService filesService
	PostgresDB   database
	MongoDB      database
}

type service struct {
	Host string
	Port string
}

type grpcHttpService struct {
	Host     string
	GrpcPort string
	HttpPort string
}

type projectsService struct {
	Service service
	DB      database
}

type authService struct {
	service
}

type logService struct {
	grpcHttpService
}

type database struct {
	User     string
	Host     string
	Port     string
	Name     string
}

type usersService struct {
	Service grpcHttpService
	DB      database
}

type postsService struct {
	grpcHttpService
}

type chatsService struct {
	Service grpcHttpService
	DB      database
}

type filesService struct {
	Protocol string
	service
}

func (cfg *Config) InitLog() {
	logSection := cfg.Section("log")
	App.LogService.Host = logSection.Key("host").String()
	App.LogService.HttpPort = logSection.Key("http_port").String()
	App.LogService.GrpcPort = logSection.Key("grpc_port").String()
}

func (cfg *Config) InitAuth() {
	authSection := cfg.Section("auth")
	App.AuthService.Host = authSection.Key("host").String()
	App.AuthService.Port = authSection.Key("port").String()
}

func (cfg *Config) InitUsers() {
	App.Users.Service = cfg.parseGrpcHttpService("users.service")
	App.Users.DB = cfg.parseDatabase("users.database")
}

func (cfg *Config) parseGrpcHttpService(sectionName string) grpcHttpService {
	section := cfg.Section(sectionName)
	return grpcHttpService{
		Host:     section.Key("host").String(),
		HttpPort: section.Key("http_port").String(),
		GrpcPort: section.Key("grpc_port").String(),
	}
}

func (cfg *Config) parseService(sectionName string) service {
	section := cfg.Section(sectionName)
	return service{
		Host: section.Key("host").String(),
		Port: section.Key("port").String(),
	}
}

func (cfg *Config) parseDatabase(sectionName string) database {
	section := cfg.Section(sectionName)
	return database{
		User:     section.Key("user").String(),
		Host:     section.Key("host").String(),
		Port:     section.Key("port").String(),
		Name:     section.Key("name").String(),
	}
}

func (cfg *Config) InitProjects() {
	App.Projects.Service = cfg.parseService("projects.service")
	App.Projects.DB = cfg.parseDatabase("projects.database")
}

func (cfg *Config) InitPosts() {
	postsSection := cfg.Section("posts")
	App.PostsService.Host = postsSection.Key("host").String()
	App.PostsService.HttpPort = postsSection.Key("http_port").String()
	App.PostsService.GrpcPort = postsSection.Key("grpc_port").String()
}

func (cfg *Config) InitChats() {
	App.Chats.Service = cfg.parseGrpcHttpService("chats.service")
	App.Chats.DB = cfg.parseDatabase("chats.database")
}

func (cfg *Config) InitFiles() {
	filesSection := cfg.Section("files")
	App.FilesService.Protocol = filesSection.Key("protocol").String()
	App.FilesService.Host = filesSection.Key("host").String()
	App.FilesService.Port = filesSection.Key("port").String()
}

func (cfg *Config) InitPostgres() {
	postgresSection := cfg.Section("postgres")
	App.PostgresDB.User = postgresSection.Key("postgres_user").String()
	App.PostgresDB.Host = postgresSection.Key("postgres_host").String()
	App.PostgresDB.Port = postgresSection.Key("postgres_port").String()
	App.PostgresDB.Name = postgresSection.Key("postgres_db_name").String()
}

func (cfg *Config) InitMongo() {
	mongoSection := cfg.Section("mongo")
	App.MongoDB.User = mongoSection.Key("mongo_user").String()
	App.MongoDB.Host = mongoSection.Key("mongo_host").String()
	App.MongoDB.Port = mongoSection.Key("mongo_port").String()
	App.MongoDB.Name = mongoSection.Key("mongo_db_name").String()
}
