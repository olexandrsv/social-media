package controller

import (
	"context"
	"fmt"
	"log"
	"social-media/database"
	"social-media/models"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	firstName := c.PostForm("firstName")
	secondName := c.PostForm("secondName")
	login := c.PostForm("login")
	passw1 := c.PostForm("passw1")
	passw2 := c.PostForm("passw2")

	if passw1 != passw2 {
		c.JSON(200, gin.H{"message": "passwords aren't equal"})
	}
	
	user, err := models.New(login, firstName, secondName, passw1)
	if err != nil{
		log.Println(err)
		c.JSON(500, "internal error")
	}

	err = user.Save()
	if err != nil{
		log.Println(err)
		c.JSON(500, "internal error")
	}

	user.Register()

	token, err := user.GenerateJWT()
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

func UserLogin(c *gin.Context) {
	login := c.PostForm("login")
	pssw := c.PostForm("passw")
	
	user := &models.User{
		Login: login,
		Password: pssw,
	}

	exist, err := user.Exist()
	if err != nil{
		log.Println(err)
		c.String(500, "internal error")
	}

	if !exist{
		log.Println(err)
		c.String(403, "forbidden")
	}

	token, err := user.GenerateJWT()
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}

	user.Register()

	c.JSON(200, gin.H{
		"token": token,
	})
}

func FollowUser(c *gin.Context) {
	followedLogin := c.Param("login")
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

	followedID, err := models.GetIdByLogin(followedLogin)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}

	user.Subscribe(followedID)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}

}

func GetFollowedInfo(c *gin.Context) {
	login := c.Param("login")
	fmt.Println(login)
	user, err := models.Load(login)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}
	c.JSON(200, user)
}

func FollowingAccounts(c *gin.Context) {
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

	followers, err := user.Followed()
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}

	var logins []string
	for _, follower := range followers{
		logins = append(logins, follower.Login)
	}
	c.JSON(200, logins)
}

func getFollowingIds(id int) ([]int, error) {
	conn := database.PostgreConn
	rows, err := conn.Query(context.Background(), "select follower_id from followers where user_id=$1", id)
	if err != nil {
		return nil, err
	}

	var following []int
	for rows.Next() {
		var user int
		err = rows.Scan(&user)
		if err != nil {
			log.Println(err)
			continue
		}
		following = append(following, user)
	}
	return following, nil
}

func getFollowingLogins(id int) ([]string, error) {
	conn := database.PostgreConn
	rows, err := conn.Query(context.Background(), "select login from users join followers on users.id = followers.user_id and followers.follower_id=$1", id)
	if err != nil {
		return nil, err
	}

	var following []string
	for rows.Next() {
		var user string
		err = rows.Scan(&user)
		if err != nil {
			log.Println(err)
			continue
		}
		following = append(following, user)
	}
	return following, nil
}
