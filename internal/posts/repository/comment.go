package repository

import (
	"context"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/posts/domain/comment"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type commentRepository interface {
	PostComments(string) ([]*comment.Comment, error)
	CreatePostComment(CreateCommentReq) (*comment.Comment, error)
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

	var comments []*comment.Comment
	for _, model := range commentModels {
		c := comment.New(model.ID, model.UserID, comment.WithText(model.Text), comment.WithImagesPaths(model.ImagesPath),
			comment.WithFilesPaths(model.FilesPath), comment.WithCommentsIDs(model.CommentsIDs))
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *repo) CreatePostComment(req CreateCommentReq) (*comment.Comment, error) {
	commentID, err := r.createComment(req)
	if err != nil {
		return nil, err
	}
	err = r.addPostChild(req.PostID, commentID)
	if err != nil {
		return nil, err
	}

	c := comment.New(commentID, req.UserID, comment.WithText(req.Text), comment.WithImagesPaths(req.ImagesPath),
		comment.WithFilesPaths(req.FilesPath))

	return c, nil
}

func (r *repo) createComment(req CreateCommentReq) (string, error) {
	model := CommentModel{
		UserID:      req.UserID,
		Text:        req.Text,
		ImagesPath:  req.ImagesPath,
		FilesPath:   req.FilesPath,
		CommentsIDs: []string{},
	}
	res, err := r.comments.InsertOne(context.Background(), model)
	if err != nil {
		log.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}
	return hexFromObjectID(res.InsertedID), nil
}
