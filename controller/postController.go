package controller

import (
	"context"
	"log"
	"mime/multipart"
	"path/filepath"
	"social-media/auth"
	"social-media/database"
	//"social-media/users"
	models "social-media/internal/users/domain/user"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostMessage(c *gin.Context) {
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
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.String(400, "invalid form param")
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
	req := generatePostRequest(text, id, imgPath, filesPath)

	coll := database.MI.DB.Collection("posts")
	result, err := coll.InsertOne(context.Background(), req)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	postId := idToHex(result.InsertedID)
	following, err := GetFollowingIds(id)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}

	req["_id"] = postId
	req["login"] = login
	req["type"] = "post"

	for _, user := range following {
		if user, ok := models.ActiveUsers.Get(user); ok {
			user.Conn.WriteJSON(req)
		}
	}

	c.JSON(200, req)
}

func processFormFiles(form *multipart.Form, c *gin.Context, key string, id int) ([]string, error) {
	filesPath := []string{}
	files := form.File[key]
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		path := strconv.Itoa(id) + "/" + filename
		filesPath = append(filesPath, path)
		if err := c.SaveUploadedFile(file, "./upload/"+path); err != nil {
			return nil, err
		}
	}
	return filesPath, nil
}

func ChangeMessage(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.String(400, "no token")
		return
	}
	login, _, err := auth.TokenCredentials(token)
	if err != nil {
		log.Println(err)
		c.String(400, "invalid credentials")
		return
	}
	rawId := c.PostForm("id")
	id, err := primitive.ObjectIDFromHex(rawId)
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
		return
	}
	text := c.PostForm("text")
	image := []string{}
	files := []string{}

	req := generatePostRequest(text, login, image, files)

	coll := database.MI.DB.Collection("posts")
	_, err = coll.UpdateByID(context.Background(), id, bson.D{{Key: "$set", Value: req}})
	if err != nil {
		log.Println(err)
		c.String(500, "internal error")
	}
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

func GetNPosts(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.String(400, "no token")
		return
	}
	id, _, err := auth.TokenCredentials(token)
	if err != nil {
		c.String(400, "invalid credentials")
		return
	}

	res, err := getPosts(id)
	if err != nil {
		c.String(500, "internal error")
		return
	}

	c.JSON(200, res)
}

func GetOtherPosts(c *gin.Context) {
	// login := c.Param("login")
	// id, err := users.GetIdByLogin(login)
	// if err != nil {
	// 	log.Println(err)
	// 	c.String(500, "internal error")
	// 	return
	// }

	// token, err := c.Cookie("token")
	// if err != nil {
	// 	log.Println(err)
	// 	c.String(400, "no token")
	// 	return
	// }
	// customerId, _, err := auth.TokenCredentials(token)
	// if err != nil {
	// 	log.Println(err)
	// 	c.String(500, "internal error")
	// 	return
	// }

	// posts, err := getPosts(id)
	// if err != nil {
	// 	log.Println(err)
	// 	c.String(500, "internal error")
	// 	return
	// }

	// num, err := countPosts("posts", "userId", id)
	// if err != nil {
	// 	c.String(500, "internal error")
	// 	return
	// }

	// conn := database.PostgreConn
	// _, err = conn.Exec(context.Background(), "update followers set read=$1 where user_id=$2 and follower_id=$3", num, id, customerId)
	// if err != nil {
	// 	log.Println(err)
	// 	c.String(500, "internal error")
	// 	return
	// }

	// c.JSON(200, posts)
}

func countPosts(collection, key string, value int) (int64, error) {
	coll := database.MI.DB.Collection(collection)
	filter := bson.M{key: value}
	count, err := coll.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getPosts(id int) ([]bson.M, error) {
	coll := database.MI.DB.Collection("posts")
	var res []bson.M
	cursor, err := coll.Find(context.Background(), bson.D{{Key: "userId", Value: id}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var post bson.M
		err := cursor.Decode(&post)
		if err != nil {
			log.Println(err)
		}
		res = append(res, post)
	}

	return res, nil
}

func idToHex(res interface{}) string {
	if res == nil {
		return ""
	}
	return res.(primitive.ObjectID).Hex()
}

func generatePostRequest(text string, id int, images, files []string) bson.M {
	req := bson.M{
		"text":     text,
		"userId":   id,
		"comments": bson.A{},
	}
	if len(images) != 0 {
		req["images"] = images
	}
	if len(files) != 0 {
		req["files"] = files
	}

	return req
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

func getIdByLogin(login string) (int, error) {
	var id int
	conn := database.PostgreConn
	err := conn.QueryRow(context.Background(), "select id from users where login=$1", login).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetFollowingIds(id int) ([]int, error) {
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
