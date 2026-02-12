package repository

import (
	"context"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/posts/domain/comment"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type commentRepository interface {
	PostComments(string) ([]*comment.Comment, error)
	GetComment(string) (*comment.Comment, error)
	CreatePostComment(CreatePostCommentReq) (*comment.Comment, error)
	DeletePostComment(string, string) error

	UpdateComment(*comment.Comment) error

	CommentComments(string) ([]*comment.Comment, error)
	CreateCommentComment(CreateCommentCommentReq) (*comment.Comment, error)
	DeleteCommentComment(string, string) error
}

func (r *repo) PostComments(postID string) ([]*comment.Comment, error) {
	commentsIDs, err := r.postCommentsIDs(postID)
	if err != nil {
		return nil, err
	}
	comments, err := r.getCommentsByIDs(commentsIDs)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *repo) getCommentsByIDs(ids []string) ([]*comment.Comment, error) {
	comments := make([]*comment.Comment, 0, len(ids))
	if len(ids) == 0 {
		return comments, nil
	}

	objectIDs, err := toObjectIDs(ids)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": objectIDs,
		},
	}
	sortOption := options.Find().SetSort(bson.M{"_id": 1})

	cur, err := r.comments.Find(context.Background(), filter, sortOption)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	var commentModels []CommentModel
	if err := cur.All(context.Background(), &commentModels); err != nil {
		return nil, common.ErrInternal
	}

	for _, model := range commentModels {
		c := comment.New(model.ID, model.UserID, model.Text, model.ImagesPath, model.FilesPath, model.CommentsIDs)
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *repo) GetComment(id string) (*comment.Comment, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}
	filter := bson.M{
		"_id": objectID,
	}
	res := r.comments.FindOne(context.Background(), filter)
	var model CommentModel
	if err := res.Decode(&model); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	return comment.New(model.ID, model.UserID, model.Text, model.ImagesPath, model.FilesPath, model.CommentsIDs), nil
}

func (r *repo) CreatePostComment(req CreatePostCommentReq) (*comment.Comment, error) {
	comment, err := r.createComment(CreateCommentReq{
		req.CreateMessageReq,
	})
	if err != nil {
		return nil, err
	}
	err = addChild(r.posts, req.PostID, comment.ID())
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *repo) createComment(req CreateCommentReq) (*comment.Comment, error) {
	model := CommentModel{
		MessageModel: MessageModel{
			UserID:     req.UserID,
			Text:       req.Text,
			ImagesPath: initIfnil(req.ImagesPaths),
			FilesPath:  initIfnil(req.FilesPaths),
		},
		CommentsIDs: []string{},
	}
	res, err := r.comments.InsertOne(context.Background(), model)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	commentID := hexFromObjectID(res.InsertedID)
	
	return comment.New(commentID, req.UserID, req.Text, req.ImagesPaths, req.FilesPaths, nil), nil
}

func (r *repo) UpdateComment(c *comment.Comment) error {
	objectID, err := primitive.ObjectIDFromHex(c.ID())
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidData
	}

	filter := bson.M{
		"_id": objectID,
	}
	update := bson.D{
		{Key: "$set", Value: UpdateCommentModel{
			UpdateMessageModel: UpdateMessageModel{
				Text:       c.Text(),
				ImagesPath: initIfnil(c.ImagesPaths()),
				FilesPath:  initIfnil(c.FilesPaths()),
			},
		}},
	}

	_, err = r.comments.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}

func (r *repo) CommentComments(commentID string) ([]*comment.Comment, error) {
	commentsIDs, err := r.commentCommentsIDs(commentID)
	if err != nil {
		return nil, err
	}

	comments, err := r.getCommentsByIDs(commentsIDs)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *repo) commentCommentsIDs(commentID string) ([]string, error) {
	objectID, err := primitive.ObjectIDFromHex(commentID)
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

	response := r.comments.FindOne(context.Background(), filter, option)

	var model CommentModel
	if err := response.Decode(&model); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	return model.CommentsIDs, nil
}

func (r *repo) CreateCommentComment(req CreateCommentCommentReq) (*comment.Comment, error) {
	comment, err := r.createComment(CreateCommentReq{
		CreateMessageReq: req.CreateMessageReq,
	})
	if err != nil {
		return nil, err
	}
	err = addChild(r.comments, req.CommentID, comment.ID())
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *repo) addCommentChild(parentID, commentID string) error {
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

	_, err = r.comments.UpdateByID(context.Background(), objectID, update)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}

func (r *repo) DeletePostComment(parentID, commentID string) error {
	session, err := r.Client.StartSession()
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		err := r.deletePostChild(sessCtx, parentID, commentID)
		if err != nil {
			if err := session.AbortTransaction(sessCtx); err != nil {
				log.Error(errors.WithStack(err))
			}
			return nil, err
		}

		err = r.deleteComment(sessCtx, commentID)
		if err != nil {
			if err := session.AbortTransaction(sessCtx); err != nil {
				log.Error(errors.WithStack(err))
			}
			return nil, err
		}

		if err := session.CommitTransaction(sessCtx); err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		return nil, nil
	})
	return err
}

func (r *repo) deletePostChild(ctx context.Context, parentID, commentID string) error {
	update := bson.M{
		"$pull": bson.M{
			"comments": commentID,
		},
	}

	objectID, err := primitive.ObjectIDFromHex(parentID)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	_, err = r.posts.UpdateByID(ctx, objectID, update)
	if err != nil {
		log.Error(errors.WithStack(err))
		return err
	}
	return nil
}

func mongoTransaction(client *mongo.Client, fns ...func()error) error{
	session, err := client.StartSession()
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		for _, fn := range fns {
			if err := fn(); err != nil {
				if err := session.AbortTransaction(sessCtx); err != nil {
					log.Error(errors.WithStack(err))
				}
				return nil, err
			}
		}
		
		if err := session.CommitTransaction(sessCtx); err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		return nil, nil
	})
	if err != nil {
		return common.ErrInternal
	}

	return nil
}

func deleteChild(ctx context.Context, coll *mongo.Collection, parentID, commentID string) error {
	update := bson.M{
		"$pull": bson.M{
			"comments": commentID,
		},
	}

	objectID, err := primitive.ObjectIDFromHex(parentID)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	_, err = coll.UpdateByID(ctx, objectID, update)
	if err != nil {
		log.Error(errors.WithStack(err))
		return err
	}
	return nil
}

func (r *repo) deleteComment(ctx context.Context, commentID string) error {
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidData
	}

	comment, err := r.GetComment(commentID)
	if err != nil {
		return err
	}
	for _, child := range comment.CommentsIDs() {
		if err = r.deleteComment(ctx, child); err != nil {
			return common.ErrInternal
		}
	}

	filter := bson.M{
		"_id": objectID,
	}
	_, err = r.comments.DeleteOne(ctx, filter)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	return nil
}

func (r *repo) DeleteCommentComment(parentID, commentID string) error {
	session, err := r.Client.StartSession()
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		err := r.deleteCommentChild(sessCtx, parentID, commentID)
		if err != nil {
			if err := session.AbortTransaction(sessCtx); err != nil {
				log.Error(errors.WithStack(err))
			}
			return nil, err
		}

		err = r.deleteComment(sessCtx, commentID)
		if err != nil {
			if err := session.AbortTransaction(sessCtx); err != nil {
				log.Error(errors.WithStack(err))
			}
			return nil, err
		}

		if err := session.CommitTransaction(sessCtx); err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		return nil, nil
	})

	return err
}

func (r *repo) deleteCommentChild(ctx context.Context, parentID, commentID string) error {
	update := bson.M{
		"$pull": bson.M{
			"comments": commentID,
		},
	}

	objectID, err := primitive.ObjectIDFromHex(parentID)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	_, err = r.posts.UpdateByID(ctx, objectID, update)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	return nil
}
