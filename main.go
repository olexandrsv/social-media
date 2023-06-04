package main

import (
	"context"
	"fmt"
	"log"
	"path"
	"path/filepath"
	"social-media/controller"
	"social-media/database"
	"social-media/middleware"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

type Config struct {
	postgresUser     string
	postgresPassword string
	postgresHost     string
	postgresPort     string
	postgresDbName   string

	mongoUser     string
	mongoPassword string
	mongoHost     string
	mongoPort     string
	mongoDbName   string
}

func main() {
	config, err := getConfig("./config/config.ini")
	if err != nil {
		log.Println(err)
		return
	}
	postgresURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.postgresUser, config.postgresPassword, config.postgresHost, config.postgresPort, config.postgresDbName)
	database.InitPostgreSQL(postgresURI)
	defer database.PostgreConn.Close()
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.mongoUser, config.mongoPassword, config.mongoHost, config.mongoPort)
	database.InitMongoDatabase(mongoURI, config.mongoDbName)
	defer database.MI.Client.Disconnect(context.Background())

	routes := gin.Default()

	routes.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./entrance/dist/entrance/index.html")
		} else {
			c.File("./entrance/dist/entrance" + path.Join(dir, file))
		}
	})

	routes.POST("/register", controller.RegisterUser)
	routes.POST("/login", controller.UserLogin)

	authorized := routes.Group("/", middleware.Auth)
	authorized.GET("/ws", controller.UpgradeToWS)
	authorized.Static("/upload", "./upload/")

	authorized.GET("/info", controller.GetUserInfo)
	authorized.PUT("/info", controller.ChangeUserInfo)
	authorized.POST("/filter", controller.GetUserByInfo)

	authorized.GET("/missed", controller.GetMissedPosts)
	authorized.GET("/missed-msg", controller.GetMissedMsg)

	authorized.POST("/follow/:login", controller.FollowUser)
	authorized.GET("/follow/:login", controller.GetFollowedInfo)

	authorized.POST("/post", controller.PostMessage)
	authorized.PUT("/post", controller.ChangeMessage)
	authorized.GET("/post", controller.GetNPosts)
	authorized.GET("/post/:login", controller.GetOtherPosts)

	authorized.GET("/follow", controller.FollowingAccounts)

	authorized.POST("/comment", controller.PostComment)
	authorized.GET("/comment/:postId", controller.GetComments)

	authorized.POST("/room", controller.NewRoom)
	authorized.GET("/rooms", controller.GetRooms)

	authorized.POST("/message", controller.ReceiveMessage)
	authorized.GET("/msg/:id", controller.GetMessages)

	routes.Run(":8080")
}

func getConfig(path string) (Config, error) {
	config := Config{}
	cfg, err := ini.Load(path)
	if err != nil {
		return config, err
	}
	pSect := cfg.Section("postgres")
	config.postgresUser = pSect.Key("postgres_user").String()
	config.postgresPassword = pSect.Key("postgres_password").String()
	config.postgresHost = pSect.Key("postgres_host").String()
	config.postgresPort = pSect.Key("postgres_port").String()
	config.postgresDbName = pSect.Key("postgres_db_name").String()

	mSect := cfg.Section("mongo")
	config.mongoUser = mSect.Key("mongo_user").String()
	config.mongoPassword = mSect.Key("mongo_password").String()
	config.mongoHost = mSect.Key("mongo_host").String()
	config.mongoPort = mSect.Key("mongo_port").String()
	config.mongoDbName = mSect.Key("mongo_db_name").String()

	return config, nil
}
