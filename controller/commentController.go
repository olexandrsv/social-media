package controller

import (
	"context"
	"log"
	"social-media/auth"
	"social-media/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PostComment(c *gin.Context) {
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
	rawId := c.PostForm("postId")
	postId, err := primitive.ObjectIDFromHex(rawId)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}
	text := c.PostForm("text")
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.String(400, "invalid param")
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

	req := generateCommentRequest(rawId, text, id, imgPath, filesPath)

	commentColl := database.MI.DB.Collection("comments")
	res, err := commentColl.InsertOne(context.Background(), req)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}
	insertedId := idToHex(res.InsertedID)
	if insertedId == "" {
		log.Println(err)
		c.String(500, "internal error")
	}
	req["login"] = login

	postsColl := database.MI.DB.Collection("posts")
	_, err = postsColl.UpdateByID(context.Background(), postId, bson.D{{"$push", bson.D{{"comments", insertedId}}}}, options.Update())
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}

	c.JSON(200, req)
}

func GetComments(c *gin.Context) {
	postId := c.Param("postId")

	coll := database.MI.DB.Collection("comments")
	var res []bson.M
	cursor, err := coll.Find(context.Background(), bson.D{{"postId", postId}})
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	conn := database.PostgreConn
	for cursor.Next(context.Background()) {
		var post bson.M
		err := cursor.Decode(&post)
		if err != nil {
			log.Println(err)
		}
		var login string
		err = conn.QueryRow(context.Background(), "select login from users where id=$1", post["id"]).Scan(&login)
		if err != nil {
			continue
		}
		post["login"] = login
		res = append(res, post)
	}

	c.JSON(200, res)
}

func generateCommentRequest(postId, text string, id int, images, files []string) bson.M {
	req := bson.M{
		"text":   text,
		"id":     id,
		"postId": postId,
	}

	if len(images) != 0 {
		req["images"] = images
	}
	if len(files) != 0 {
		req["files"] = files
	}
	return req
}
