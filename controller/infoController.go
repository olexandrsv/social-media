package controller

import (
	"context"
	"log"
	"social-media/auth"
	"social-media/database"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.String(400, "no token")
		return
	}
	id, login, err := auth.TokenCredentials(token)
	if err != nil {
		log.Println(err)
		c.String(400, "invalid credentials")
		return
	}
	conn := database.PostgreConn
	var firstName string
	var secondName string
	var bio string
	var interests string
	err = conn.QueryRow(context.Background(), "select first_name, second_name, bio, interests from users where id=$1", id).Scan(&firstName, &secondName, &bio, &interests)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}
	c.JSON(200, gin.H{
		"login":      login,
		"firstName":  firstName,
		"secondName": secondName,
		"bio":        bio,
		"interests":  interests,
	})
}

func ChangeUserInfo(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.String(400, "no token")
		return
	}
	id, _, err := auth.TokenCredentials(token)
	if err != nil {
		log.Println(err)
		c.String(400, "invalid credentials")
		return
	}
	firstName := c.PostForm("firstName")
	secondName := c.PostForm("secondName")
	bio := c.PostForm("bio")
	interests := c.PostForm("interests")

	conn := database.PostgreConn
	_, err = conn.Exec(context.Background(), "update users set first_name=$1, second_name=$2, bio=$3, interests=$4 where id=$5", firstName, secondName, bio, interests, id)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}
}

func GetUserByInfo(c *gin.Context) {
	word := c.PostForm("word")

	conn := database.PostgreConn
	rows, err := conn.Query(context.Background(), "select login from users where interests like $1 or bio like $1", "%"+word+"%")
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}

	var sb strings.Builder
	for rows.Next() {
		var login string
		err = rows.Scan(&login)
		if err != nil {
			continue
		}
		sb.WriteString(",")
		sb.WriteString(login)
	}
	res := sb.String()
	if len(res) > 0 {
		res = res[1:]
	}
	c.JSON(200, res)
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
	id, _, err := auth.TokenCredentials(token)
	if err != nil {
		log.Println(err)
		c.String(400, "invalid credentials")
		return
	}
	following, err := getFollowingLogins(id)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	res := make(map[string]int)

	for _, user := range following {
		userId, err := getIdByLogin(user)
		if err != nil {
			log.Println(err)
			continue
		}
		count, err := countPosts("posts", "userId", userId)
		if err != nil {
			log.Println(err)
			continue
		}
		conn := database.PostgreConn
		var readNum int
		row := conn.QueryRow(context.Background(), "select read from followers where user_id=$1 and follower_id=$2", userId, id)
		if err := row.Scan(&readNum); err != nil {
			log.Println(err)
			continue
		}
		res[user] = int(count) - readNum
	}

	log.Println(res)

	c.JSON(200, res)
}
