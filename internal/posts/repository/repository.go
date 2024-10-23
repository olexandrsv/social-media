package repository

import (
	"context"
	"fmt"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/posts/domain/post"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	CreatePost(PostModel) (*post.Post, error)
}

type repo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func New() Repository {
	user := config.App.MongoDB.User
	password := config.App.MongoDB.Password
	host := config.App.MongoDB.Host
	port := config.App.MongoDB.Port
	databaseName := config.App.MongoDB.Name

	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", user, password, host, port, databaseName)
	fmt.Println("url: ", url)
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	db := client.Database(databaseName)

	return &repo{
		Client: client,
		DB:     db,
	}
}

func (r *repo) CreatePost(postModel PostModel) (*post.Post, error) {
	coll := r.DB.Collection("mock")
	res, err := coll.InsertOne(context.Background(), postModel)
	if err != nil {
		log.Error(err)
		return nil, common.ErrInternal
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()

	post := post.New(id, postModel.UserID, post.WithText(postModel.Text),
		post.WithFilesPaths(postModel.FilesPath), post.WithImagesPaths(postModel.ImagesPath))
	return post, nil
}
