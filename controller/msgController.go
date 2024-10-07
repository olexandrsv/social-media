package controller

import (
	"context"
	"log"
	"social-media/auth"
	"social-media/database"
	models "social-media/internal/users/domain/user"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func ReceiveMessage(c *gin.Context) {
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
	text := c.PostForm("text")
	roomId, err := strconv.Atoi(c.PostForm("roomId"))
	if err != nil {
		c.String(400, "invalid form field")
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.String(400, "invalid form fields")
		return
	}

	imgPath, err := processFormFiles(form, c, "images[]", id)
	if err != nil {
		log.Println(err)
		c.String(400, "upload file err")
		return
	}
	filesPath, err := processFormFiles(form, c, "files[]", id)
	if err != nil {
		log.Println(err)
		c.String(400, "upload file err")
		return
	}

	req := generateMsgRequest(text, id, roomId, imgPath, filesPath)

	coll := database.MI.DB.Collection("messages")
	_, err = coll.InsertOne(context.Background(), req)
	if err != nil {
		c.String(500, "internal error")
	}

	following, err := getRoomUsers(roomId)
	if err != nil {
		return
	}

	req["login"] = login
	req["type"] = "msg"

	for _, user := range following {
		if user == id {
			continue
		}
		if user, ok := models.ActiveUsers.Get(user); ok {
			user.Conn.WriteJSON(req)
		}
	}

	c.JSON(200, req)
}

func getRoomUsers(roomId int) ([]int, error) {
	conn := database.PostgreConn
	rows, err := conn.Query(context.Background(), "select user_id from urooms where room_id=$1", strconv.Itoa(roomId))
	if err != nil {
		return nil, err
	}

	var users []int
	for rows.Next() {
		var userId int
		err = rows.Scan(&userId)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, userId)
	}
	return users, nil
}

func GetMessages(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.String(400, "no token")
		return
	}
	userId, _, err := auth.TokenCredentials(token)
	if err != nil {
		log.Println(err)
		c.String(400, "invalid credentials")
		return
	}

	roomId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.String(400, "invalid param")
		return
	}

	res, err := getMsgs(roomId)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	count, err := countPosts("messages", "roomId", roomId)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	conn := database.PostgreConn
	_, err = conn.Exec(context.Background(), "update read_msg set count=$1 where user_id=$2 and room_id=$3", count, userId, roomId)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	for _, msg := range res {
		var login string
		err := conn.QueryRow(context.Background(), "select login from users where id=$1", msg["userId"]).Scan(&login)
		if err != nil {
			continue
		}
		msg["login"] = login
	}

	c.JSON(200, res)
}

func GetMissedMsg(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.String(400, "no token")
		return
	}
	id, _, err := auth.TokenCredentials(token)
	if err != nil {
		log.Println(err)
		c.String(400, "invalid creadentials")
		return
	}

	var rooms []int

	conn := database.PostgreConn
	rows, err := conn.Query(context.Background(), "select room_id from urooms where user_id=$1", id)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}

	for rows.Next() {
		var roomId int
		err := rows.Scan(&roomId)
		if err != nil {
			log.Println(err)
			continue
		}
		rooms = append(rooms, roomId)
	}

	res := make(map[int]int)

	for _, roomId := range rooms {
		if err != nil {
			log.Println(err)
			continue
		}
		count, err := countPosts("messages", "roomId", roomId)
		if err != nil {
			log.Println(err)
			continue
		}
		conn := database.PostgreConn
		var readNum int
		row := conn.QueryRow(context.Background(), "select count from read_msg where user_id=$1 and room_id=$2", id, roomId)
		if err := row.Scan(&readNum); err != nil {
			log.Println(err)
			continue
		}
		res[roomId] = int(count) - readNum
	}

	c.JSON(200, res)
}

func getRoomById(id int) (string, error) {
	var name string
	conn := database.PostgreConn
	err := conn.QueryRow(context.Background(), "select name from rooms where id=$1", id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func generateMsgRequest(text string, id, roomId int, images, files []string) bson.M {
	req := bson.M{
		"text":   text,
		"userId": id,
		"roomId": roomId,
	}
	if len(images) != 0 {
		req["images"] = images
	}
	if len(files) != 0 {
		req["files"] = files
	}

	return req
}

func getMsgs(id int) ([]bson.M, error) {
	coll := database.MI.DB.Collection("messages")
	var res []bson.M
	cursor, err := coll.Find(context.Background(), bson.D{{Key: "roomId", Value: id}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var post bson.M
		err := cursor.Decode(&post)
		if err != nil {
			log.Println(err)
			continue
		}
		res = append(res, post)
	}

	return res, nil
}
