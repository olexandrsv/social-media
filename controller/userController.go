package controller

import (
	"context"
	"log"
	"social-media/auth"
	"social-media/database"
	"social-media/hash"
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
	hashPsw, err := hash.HashPassword(passw1)
	if err != nil {
		c.JSON(500, gin.H{})
	}

	var id int

	conn := database.PostgreConn
	err = conn.QueryRow(context.Background(), "insert into users (login, first_name, second_name, password, bio, interests) values ($1, $2, $3, $4, $5, $6) returning id", login, firstName, secondName, hashPsw, "", "").Scan(&id)
	if err != nil {
		log.Println(err)
		c.JSON(500, "internal error")
		return
	}

	token, err := auth.GenerateJWT(id, login)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	user := &models.User{
		Id:    id,
		Login: login,
	}
	models.ActiveUsers.Set(id, user)

	c.JSON(200, gin.H{
		"token": token,
	})
}

func UserLogin(c *gin.Context) {
	login := c.PostForm("login")
	pssw := c.PostForm("passw")
	var id int
	var encodedPassw string

	conn := database.PostgreConn
	err := conn.QueryRow(context.Background(), "select id, password from users where login=$1", login).Scan(&id, &encodedPassw)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	if !hash.CheckPassword(pssw, encodedPassw) {
		log.Println(err)
		c.String(403, "forbidden")
		return
	}

	token, err := auth.GenerateJWT(id, login)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	user := &models.User{
		Id:    id,
		Login: login,
	}
	models.ActiveUsers.Set(id, user)

	c.JSON(200, gin.H{
		"token": token,
	})
}

func FollowUser(c *gin.Context) {
	userLogin := c.Param("login")
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

	followerId, err := getIdByLogin(userLogin)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	conn := database.PostgreConn
	_, err = conn.Exec(context.Background(), "insert into followers (user_id, follower_id) values ($1, $2)", followerId, id)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}

}

func getIdByLogin(login string) (int, error) {
	var id int
	conn := database.PostgreConn
	err := conn.QueryRow(context.Background(), "select id from users where login=$1", login).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetFollowedInfo(c *gin.Context) {
	login := c.Param("login")
	var firstName string
	var secondName string
	var bio string
	var interests string
	conn := database.PostgreConn
	err := conn.QueryRow(context.Background(), "select first_name, second_name, bio, interests from users where login=$1", login).Scan(&firstName, &secondName, &bio, &interests)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}
	c.JSON(200, gin.H{
		"login":      login,
		"firstName":  firstName,
		"secondName": secondName,
		"bio":        bio,
		"interests":  interests,
	})
}

func FollowingAccounts(c *gin.Context) {
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

	c.JSON(200, following)
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
