package repository

import (
	"context"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"
	"social-media/internal/posts/domain/chatmessage"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type chatMessageRepository interface {
	ChatMessages(int) ([]*chatmessage.ChatMessage, error)
	ChatMessage(string) (*chatmessage.ChatMessage, error)
	CreateChatMessage(CreateChatMessageReq) (*chatmessage.ChatMessage, error)
	UpdateChatMessage(*chatmessage.ChatMessage) error
	DeleteChatMessage(string) error
	MissedMessagesNumber(int, []*chatmessage.ChatMessage) ([]*chatmessage.MissedNumber, error)
}

func (r *repo) ChatMessages(chatID int) ([]*chatmessage.ChatMessage, error) {
	filter := bson.M{
		"chatId": chatID,
	}

	cursor, err := r.messages.Find(context.Background(), filter)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	var chatMessages []ChatMessageModel
	if err := cursor.All(context.Background(), &chatMessages); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	chats := slice.MustConvert(chatMessages, func(m ChatMessageModel) *chatmessage.ChatMessage {
		return chatmessage.New(m.ID, m.UserID, m.Text, m.ImagesPath, m.FilesPath, m.ChatID)
	})

	return chats, nil
}

func (r *repo) ChatMessage(id string) (*chatmessage.ChatMessage, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	filter := bson.M{
		"_id": objectID,
	}

	res := r.messages.FindOne(context.Background(), filter)
	var m ChatMessageModel
	if err := res.Decode(&m); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	return chatmessage.New(m.ID, m.UserID, m.Text, m.ImagesPath, m.FilesPath, m.ChatID), nil
}

func (r *repo) CreateChatMessage(req CreateChatMessageReq) (*chatmessage.ChatMessage, error) {
	model := ChatMessageModel{
		MessageModel: MessageModel{
			UserID:     req.UserID,
			Text:       req.Text,
			ImagesPath: initIfnil(req.ImagesPaths),
			FilesPath:  initIfnil(req.FilesPaths),
		},
		ChatID: req.ChatID,
	}
	res, err := r.messages.InsertOne(context.Background(), model)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()
	return chatmessage.New(id, req.UserID, req.Text, req.ImagesPaths, req.FilesPaths, req.ChatID), nil
}

func (r *repo) UpdateChatMessage(m *chatmessage.ChatMessage) error {
	objectID, err := primitive.ObjectIDFromHex(m.ID())
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidData
	}

	filter := bson.M{
		"_id": objectID,
	}
	update := bson.D{
		{Key: "$set", Value: UpdateChatMessageModel{
			UpdateMessageModel: UpdateMessageModel{
				Text:       m.Text(),
				ImagesPath: initIfnil(m.ImagesPaths()),
				FilesPath:  initIfnil(m.FilesPaths()),
			},
		}},
	}

	_, err = r.messages.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	return nil
}

func (r *repo) DeleteChatMessage(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidData
	}

	filter := bson.M{
		"_id": objectID,
	}
	_, err = r.messages.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}

	return nil
}

func (r *repo) MissedMessagesNumber(userID int, lastReadMessages []*chatmessage.ChatMessage) ([]*chatmessage.MissedNumber, error) {
	if len(lastReadMessages) == 0 {
		return make([]*chatmessage.MissedNumber, 0), nil
	}
	for _, m := range lastReadMessages{
		log.Infof("Message{ ID: %s, ChatID: %d }", m.ID(), m.ChatID())
	}
	var orConditions []bson.M
	for _, lastReadMessage := range lastReadMessages {
		var objectID primitive.ObjectID
		if lastReadMessage.ID() == "" {
			objectID = primitive.NilObjectID
		} else {
			id, err := primitive.ObjectIDFromHex(lastReadMessage.ID())
			if err != nil {
				log.Error(errors.WithStack(err))
				return nil, common.ErrInternal
			}
			objectID = id
		}
		orConditions = append(orConditions, bson.M{
			"chatId": lastReadMessage.ChatID(),
			"_id":    bson.M{"$gt": objectID},
			"userId": bson.M{ "$ne": userID},
		})
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"$or": orConditions}}},
		{{"$group", bson.M{
			"_id":    "$chatId",
			"number": bson.M{"$sum": 1},
			"ids":    bson.M{"$push": "$_id"},
		}}},
	}

	cursor, err := r.messages.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	var missedMessagesNumbers []*chatmessage.MissedNumber
	for cursor.Next(context.Background()) {
		var missedMessagesNumber MissedMessagesNumber
		if err := cursor.Decode(&missedMessagesNumber); err != nil {
			log.Error(errors.WithStack(err))
			return nil, common.ErrInternal
		}
		missedMessageNumber := chatmessage.NewMissedNumber(missedMessagesNumber.ChatID, missedMessagesNumber.Number)
		missedMessagesNumbers = append(missedMessagesNumbers, missedMessageNumber)
	}

	return missedMessagesNumbers, nil
}
