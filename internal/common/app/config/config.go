package config

var App AppConfig

type AppConfig struct{
	UsersService usersService
	AuthService  authService
	LogService   logService
	Database     database
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
