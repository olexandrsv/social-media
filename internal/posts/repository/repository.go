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
	CreatePost(CreatePostReq) (*post.Post, error)
	UserPosts(int) ([]*post.Post, error)
	UpdatePost(*post.Post) error
	GetPost(string) (*post.Post, error)
	DeletePost(string) error
	MissedPostsNumber([]*post.LastRead) ([]*post.MissedNumber, error)
	commentRepository
	chatMessageRepository
}

type repo struct {
	Client       *mongo.Client
	DB           *mongo.Database
	posts        *mongo.Collection
	comments     *mongo.Collection
	messages     *mongo.Collection
	readMessages *mongo.Collection
}

func New() (Repository, error) {
	password, err := common.ReadSecret("mongodb-secret")
	if err != nil {
		return nil, err
	}
	user := config.App.MongoDB.User
	host := config.App.MongoDB.Host
	port := config.App.MongoDB.Port
	databaseName := config.App.MongoDB.Name

	// mongodb://mongo_adiutor:27017/?replicaSet=rs1&directConnection=true
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?replicaSet=rs1&authSource=admin", user, password, host, port, databaseName)
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
		Client:       client,
		DB:           db,
		posts:        db.Collection("posts"),
		comments:     db.Collection("comments"),
		messages:     db.Collection("messages"),
		readMessages: db.Collection("read_messages"),
	}, nil
}

func (r *repo) CreatePost(req CreatePostReq) (*post.Post, error) {
	model := PostModel{
		MessageModel: MessageModel{
			UserID:     req.UserID,
			Text:       req.Text,
			ImagesPath: initIfnil(req.ImagesPaths),
			FilesPath:  initIfnil(req.FilesPaths),
		},
		CommentsIDs: []string{},
	}
	res, err := r.posts.InsertOne(context.Background(), model)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()

	return post.New(id, model.UserID, model.Text, model.ImagesPath, model.FilesPath, nil), nil
}

func (r *repo) UserPosts(userID int) ([]*post.Post, error) {
	coll := r.DB.Collection("posts")
	filter := bson.D{{Key: "userId", Value: userID}}
	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	var postModels []PostModel
	if err := cursor.All(context.Background(), &postModels); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	posts := make([]*post.Post, 0, len(postModels))
	for _, postModel := range postModels {
		post := post.New(postModel.ID, postModel.UserID, postModel.Text, postModel.ImagesPath, postModel.FilesPath, nil)
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *repo) UpdatePost(post *post.Post) error {
	coll := r.DB.Collection("posts")
	id, err := primitive.ObjectIDFromHex(post.ID())
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidData
	}

	update := bson.D{
		{Key: "$set", Value: UpdatePostModel{
			UpdateMessageModel: UpdateMessageModel{
				Text:       post.Text(),
				ImagesPath: initIfnil(post.ImagesPaths()),
				FilesPath:  initIfnil(post.FilesPaths()),
			},
		}},
	}
	_, err = coll.UpdateByID(context.Background(), id, update)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}

func initIfnil(slice []string) []string {
	if slice == nil {
		return []string{}
	}
	return slice
}

func (r *repo) GetPost(id string) (*post.Post, error) {
	coll := r.DB.Collection("posts")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	filter := bson.M{
		"_id": objectID,
	}

	res := coll.FindOne(context.Background(), filter)
	var model PostModel
	err = res.Decode(&model)
	if err == mongo.ErrNoDocuments {
		return nil, common.ErrNotFound
	}
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	return post.New(model.ID, model.UserID, model.Text, model.ImagesPath, model.FilesPath, model.CommentsIDs), nil
}

func (r *repo) DeletePost(id string) error {
	session, err := r.Client.StartSession()
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		err := r.deletePost(sessCtx, id)
		if err != nil {
			if err := session.AbortTransaction(sessCtx); err != nil {
				log.Error(errors.WithStack(err))
			}
			return nil, err
		}

		if err = session.CommitTransaction(sessCtx); err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}

		return nil, nil
	})

	return err
}

func (r *repo) deletePost(ctx context.Context, id string) error {
	post, err := r.GetPost(id)
	if err != nil {
		return err
	}
	for _, commentID := range post.CommentsIDs() {
		if err := r.deleteComment(ctx, commentID); err != nil {
			return err
		}
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidData
	}

	filter := bson.M{
		"_id": objectID,
	}
	_, err = r.posts.DeleteOne(ctx, filter)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}

func (r *repo) postCommentsIDs(postID string) ([]string, error) {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	filter := bson.M{
		"_id": objectID,
	}
	option := options.FindOne().SetProjection(bson.M{
		"comments": 1,
	})

	response := r.posts.FindOne(context.Background(), filter, option)

	var result PostModel
	if err := response.Decode(&result); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	return result.CommentsIDs, nil
}

func toObjectIDs(ids []string) ([]primitive.ObjectID, error) {
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		objectIDs = append(objectIDs, objectID)
	}
	return objectIDs, nil
}

func (r *repo) addPostChild(postID, commentID string) error {
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidData
	}

	update := bson.M{
		"$push": bson.M{
			"comments": commentID,
		},
	}

	_, err = r.posts.UpdateByID(context.Background(), objectID, update)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}

func addChild(coll *mongo.Collection, parentID, commentID string) error {
	objectID, err := primitive.ObjectIDFromHex(parentID)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidData
	}

	update := bson.M{
		"$push": bson.M{
			"comments": commentID,
		},
	}

	_, err = coll.UpdateByID(context.Background(), objectID, update)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}

func hexFromObjectID(i interface{}) string {
	return i.(primitive.ObjectID).Hex()
}

func (r *repo) MissedPostsNumber(lastReadPosts []*post.LastRead) ([]*post.MissedNumber, error) {
	if len(lastReadPosts) == 0 {
		return make([]*post.MissedNumber, 0), nil
	}
	var orConditions []bson.M
	for _, item := range lastReadPosts {
		var objectID primitive.ObjectID
		if item.LastReadPostID() == "" {
			objectID = primitive.NilObjectID
		} else {
			id, err := primitive.ObjectIDFromHex(item.LastReadPostID())
			if err != nil {
				log.Error(errors.WithStack(err))
				return nil, common.ErrInternal
			}
			objectID = id
		}
		orConditions = append(orConditions, bson.M{
			"userId": item.FollowingID(),
			"_id":    bson.M{"$gt": objectID},
		})
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"$or": orConditions}}},
		{{"$group", bson.M{
			"_id":    "$userId",
			"number": bson.M{"$sum": 1},
			"ids":    bson.M{"$push": "$_id"},
		}}},
	}

	cursor, err := r.posts.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	var missedPostsNumbers []*post.MissedNumber
	for cursor.Next(context.Background()) {
		var missedPostsNumber MissedPostsNumber
		if err := cursor.Decode(&missedPostsNumber); err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		missedPostNumber := post.NewMissedNumber(missedPostsNumber.FollowingID, missedPostsNumber.Number)
		log.Infof("missedPost{ followingID: %d, number: %d }", missedPostsNumber.FollowingID, missedPostsNumber.Number)
		missedPostsNumbers = append(missedPostsNumbers, missedPostNumber)
	}

	return missedPostsNumbers, nil
}
