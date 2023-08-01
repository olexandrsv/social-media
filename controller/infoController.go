package controller

import (
	"context"
	"log"
	"social-media/auth"
	"social-media/database"
	"social-media/models"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.String(400, "no token")
	}
	user, err := models.ParseToken(token)
	if err != nil {
		log.Println(err)
		c.String(400, "invalid credentials")
	}

	user, err = models.Load(user.Login)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}
	c.JSON(200, user)
}

func ChangeUserInfo(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.String(400, "no token")
	}

	user, err := models.ParseToken(token)
	if err != nil {
		log.Println(err)
		c.String(400, "invalid credentials")
	}

	user.FirstName = c.PostForm("firstName")
	user.SecondName = c.PostForm("secondName")
	user.Bio = c.PostForm("bio")
	user.Interests = c.PostForm("interests")

	if err = user.Update(); err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}
}

func GetUserByInfo(c *gin.Context) {
	word := c.PostForm("word")

	logins, err := models.GetByInfo(word)
	if err != nil {
		log.Println(err)
		c.JSON(500, "internal error")
	}
	c.JSON(200, logins)
}

func getLogin(c *gin.Context) (string, error) {
	token, err := c.Cookie("token")
	if err != nil {
		return "", err
	}
	_, login, err := auth.TokenCredentials(token)
	if err != nil {
		c.String(500, err.Error())
		return "", err
	}
	return login, nil
}

func GetMissedPosts(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.String(400, "no token")
		return
	}

	user, err := models.ParseToken(token)
	if err != nil {
		log.Println(err)
		c.String(400, "invalid credentials")
		return
	}

	users, err := user.Followed()
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	res := make(map[string]int)

	for _, followed := range users {
		if err != nil {
			log.Println(err)
			continue
		}
		count, err := countPosts("posts", "userId", followed.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		conn := database.PostgreConn
		var readNum int
		row := conn.QueryRow(context.Background(), "select read from followers where user_id=$1 and follower_id=$2", followed.ID, user.ID)
		if err := row.Scan(&readNum); err != nil {
			log.Println(err)
			continue
		}
		res[user.Login] = int(count) - readNum
	}

	log.Println(res)

	c.JSON(200, res)
}
