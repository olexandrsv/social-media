package config

import "gopkg.in/ini.v1"

type usersService struct {
	Port string
}

func InitUsersConfig() {
	//cfg, err := ini.Load("./../../internal/common/app/config/config.ini")
	cfg, err := ini.Load("./../../common/app/config/config.ini")
	if err != nil {
		panic(err)
	}

	logSection := cfg.Section("log")
	App.LogService.Port = logSection.Key("port").String()

	authSection := cfg.Section("auth")
	App.AuthService.Port = authSection.Key("port").String()

	usersSection := cfg.Section("users")
	App.UsersService.Port = usersSection.Key("port").String()

	postgresSection := cfg.Section("postgres")
	App.Database.User = postgresSection.Key("postgres_user").String()
	App.Database.Password = postgresSection.Key("postgres_password").String()
	App.Database.Host = postgresSection.Key("postgres_host").String()
	App.Database.Port = postgresSection.Key("postgres_port").String()
	App.Database.Name = postgresSection.Key("postgres_db_name").String()
}
