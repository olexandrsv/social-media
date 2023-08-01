package controller

import (
	"context"
	"log"
	"net/http"
	"social-media/auth"
	"social-media/database"
	"social-media/models"
	"social-media/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func UpgradeToWS(c *gin.Context) {
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
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}
	defer conn.Close()
	if user, ok := models.ActiveUsers.Get(id); ok {
		user.Conn = conn
	} else {
		newUsers := &models.User{
			ID:    id,
			Login: login,
			Conn:  conn,
		}
		models.ActiveUsers.Set(id, newUsers)
	}

	var ids []int
	psql := database.PostgreConn
	rows, err := psql.Query(context.Background(), "select id from rooms join urooms on rooms.id=urooms.room_id and urooms.user_id=$1", id)
	if err != nil {
		log.Println(err)
		c.String(500, err.Error())
		return
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	for _, id := range ids {
		if _, ok := models.ActiveRoom.Get(id); !ok {
			var name string
			var users []int
			err = psql.QueryRow(context.Background(), "select name, array_agg(user_id) from rooms join urooms on rooms.id=urooms.room_id and rooms.id=$1 group by rooms.id, urooms.name", id).Scan(&name, &users)
			if err != nil {
				log.Println(err)
				continue
			}
			room := &models.Room{
				Id:    id,
				Name:  name,
				Users: users,
			}
			models.ActiveRoom.Set(id, room)
		}
	}

	ws.WSHandler(conn)
}
