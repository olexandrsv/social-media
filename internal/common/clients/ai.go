package clients

import (
	"bytes"
	"encoding/json"
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"
	"social-media/internal/posts/domain/message"
	"social-media/internal/posts/domain/tone"

	"github.com/pkg/errors"
)

type AIClient[T message.MessageI] interface {
	EstimateTone([]T) (*tone.Tone, error)
}

type aiClient[T message.MessageI] struct {
}

func NewAIClient[T message.MessageI]() AIClient[T] {
	return &aiClient[T]{}
}

type EstimateToneReq struct {
	Messages []MessageModel `json:"messages"`
}

type MessageModel struct {
	Text string `json:"text"`
}

type EstimateToneResp struct {
	Tone Tone `json:"tone_estimation"`
}

type Tone struct {
	Positive float64 `json:"positive_percentage"`
	Negative float64 `json:"negative_percentage"`
}

func (c *aiClient[T]) EstimateTone(messages []T) (*tone.Tone, error) {
	messagesModels := slice.MustConvert(messages, func(message T) MessageModel {
		return MessageModel{
			Text: message.Text(),
		}
	})

	data, err := json.Marshal(EstimateToneReq{
		Messages: messagesModels,
	})
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	req, err := http.NewRequest("POST", "http://ai:8000/data", bytes.NewBuffer(data))
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	parsedResp, err := common.ParseResponse[EstimateToneResp](resp)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	return tone.New(parsedResp.Tone.Positive, parsedResp.Tone.Negative), nil
}
