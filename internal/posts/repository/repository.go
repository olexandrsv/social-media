package repository

import (
	"context"
	"fmt"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/posts/domain/post"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	CreatePost(PostModel) (*post.Post, error)
	UserPosts(int) ([]*post.Post, error)
	UpdatePost(*post.Post) error
	GetPost(string) (*post.Post, error)
	DeletePost(string) error
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
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Error(errors.WithStack(err))
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Error(errors.WithStack(err))
		panic(err)
	}
	db := client.Database(databaseName)

	return &repo{
		Client: client,
		DB:     db,
	}
}

func (r *repo) CreatePost(postModel PostModel) (*post.Post, error) {
	coll := r.DB.Collection("posts")
	res, err := coll.InsertOne(context.Background(), postModel)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()

	post := post.New(id, postModel.UserID, post.WithText(postModel.Text),
		post.WithFilesPaths(postModel.FilesPath), post.WithImagesPaths(postModel.ImagesPath))
	return post, nil
}

func (r *repo) UserPosts(userID int) ([]*post.Post, error){
	coll := r.DB.Collection("posts")
	filter := bson.D{{Key: "userId", Value: userID}}
	cursor, err := coll.Find(context.Background(), filter)
	if err != nil{
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	var postModels []PostModel
	if err := cursor.All(context.Background(), &postModels); err != nil{
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	posts := make([]*post.Post, 0, len(postModels))
	for _, postModel := range postModels{
		post := post.New(postModel.ID, postModel.UserID, post.WithText(postModel.Text),
			post.WithImagesPaths(postModel.ImagesPath), post.WithFilesPaths(postModel.FilesPath))
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *repo) UpdatePost(post *post.Post) error{
	coll := r.DB.Collection("posts")
	id, err := primitive.ObjectIDFromHex(post.ID())
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidData
	}

	update := bson.D{
		{Key: "$set", Value: UpdatePostModel{
			Text: post.Text(),
			ImagesPath: post.ImagesPaths(),
			FilesPath: post.FilesPaths(),
		}},
	}
	_, err = coll.UpdateByID(context.Background(), id, update)
	if err != nil{
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}

func (r *repo) GetPost(id string) (*post.Post, error){
	coll := r.DB.Collection("posts")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	filter := bson.M{
		"_id": objectID,
	}
	
	res := coll.FindOne(context.Background(), filter)
	var model PostModel
	err = res.Decode(&model)
	if err == mongo.ErrNoDocuments{
		return nil, common.ErrNotFound
	}
	if err != nil{
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	
	return post.New(model.ID, model.UserID, post.WithText(model.Text), post.WithImagesPaths(model.ImagesPath),
		post.WithFilesPaths(model.FilesPath), post.WithCommentsIDs(model.CommentsIDs)), nil
}

func (r *repo) DeletePost(id string) error {
	coll := r.DB.Collection("posts")
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		log.Error(errors.WithStack(err))
		return common.ErrInvalidData
	}

	filter := bson.M{
		"_id": objectID,
	}
	_, err = coll.DeleteOne(context.Background(), filter)
	if err != nil{
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}
