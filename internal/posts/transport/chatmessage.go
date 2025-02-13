package transport

import (
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"
	"social-media/internal/posts/service"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) getChatMessages(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrInvalidToken)
		return
	}

	rawID, ok := mux.Vars(r)["chat_id"]
	if !ok {
		writeError(w, common.ErrInvalidData)
		return
	}
	id, err := strconv.Atoi(rawID)
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrInvalidData)
		return
	}

	messages, err := s.service.ChatMessages(token.Value, id)
	if err != nil {
		writeError(w, err)
		return
	}

	models := slice.MustConvert(messages, chatMessageToModel)
	writeJSON(w, models)
}

func (s *server) createChatMessage(w http.ResponseWriter, r *http.Request) {
	req, err := s.decodeCreateMessageReq(r)
	if err != nil {
		writeError(w, err)
		return
	}

	rawID, ok := mux.Vars(r)["chat_id"]
	if !ok {
		writeError(w, common.ErrInvalidData)
		return
	}
	id, err := strconv.Atoi(rawID)
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrInvalidData)
		return
	}

	m, err := s.service.CreateChatMessage(service.CreateChatMessageReq{
		ChatID:           id,
		CreateMessageReq: req,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, chatMessageToModel(m))
}

func (s *server) updateChatMessage(w http.ResponseWriter, r *http.Request) {
	req, err := s.decodeUpdateMessageReq(r)
	if err != nil {
		writeError(w, err)
		return
	}

	err = s.service.UpdateChatMessage(service.UpdateChatMessageReq{
		UpdateMessageReq: req,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, nil)
}

func (s *server) deleteChatMessage(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrNoToken)
		return
	}

	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrInvalidData)
		return
	}

	err = s.service.DeleteChatMessage(token.Value, id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, nil)
}
