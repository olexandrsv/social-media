package controller

import (
	"context"
	"log"
	"social-media/auth"
	"social-media/database"
	"social-media/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewRoom(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.String(400, "no token")
		return
	}
	_, login, err := auth.TokenCredentials(token)
	if err != nil {
		c.String(400, "invalid credentials")
		return
	}

	name := c.PostForm("name")
	list := c.PostForm("users")
	users := strings.Split(list, ", ")
	users = append(users, login)

	var roomId int

	conn := database.PostgreConn
	conn.QueryRow(context.Background(), "insert into rooms (name) values ($1) returning id", name).Scan(&roomId)

	var userIds []int
	for _, user := range users {
		userId, err := getIdByLogin(user)
		if err != nil {
			continue
		}
		conn.Exec(context.Background(), "insert into urooms (room_id, user_id) values ($1, $2)", roomId, userId)
		conn.Exec(context.Background(), "insert into read_msg (room_id, user_id, count) values ($1, $2, $3)", roomId, userId, 0)
		userIds = append(userIds, userId)
	}

	room := &models.Room{
		Id:    roomId,
		Name:  name,
		Users: userIds,
	}
	models.ActiveRoom.Set(roomId, room)

	c.JSON(200, room)
}

func GetRooms(c *gin.Context) {
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
	conn := database.PostgreConn
	rows, err := conn.Query(context.Background(), "select id, name from rooms join urooms on rooms.id=urooms.room_id where user_id=$1", id)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		if err = rows.Scan(&room.Id, &room.Name); err != nil {
			continue
		}
		rooms = append(rooms, room)
	}

	c.JSON(200, rooms)
}
