package transport

import (
	"fmt"
	"io"
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"
	"social-media/internal/posts/domain/chatmessage"
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
		log.Error(errors.New("chat_id doesn't exists"))
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
		log.Error(errors.New("chat_id doesn't exists"))
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
	rawChatID := r.FormValue("chat_id")

	chatID, err := strconv.Atoi(rawChatID)
	if err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	m, err := s.service.UpdateChatMessage(service.UpdateChatMessageReq{
		ChatID:           chatID,
		UpdateMessageReq: req,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, chatMessageToModel(m))
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

func (s *server) getMissedMessages(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrInvalidToken)
		return
	}

	messages, err := s.service.MissedMessagesNumber(token.Value)
	if err != nil {
		writeError(w, err)
		return
	}

	convertedMessages := slice.MustConvert(messages, func(m *chatmessage.MissedNumber) MissedMessagesNumber {
		return MissedMessagesNumber{
			ChatID: m.ChatID(),
			Number: m.Number(),
		}
	})

	writeJSON(w, convertedMessages)
}

func (s *server) getChatMessageFile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	signedUrl, ok := mux.Vars(r)["signed_url"]
	if !ok {
		writeError(w, common.ErrInvalidData)
		return
	}

	file, err := s.service.GetChatMessageFile(cookie.Value, signedUrl)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name()))
	w.Header().Set("Content-Type", "text/plain")

	_, err = io.Copy(w, file)
	if err != nil {
		log.Error(errors.WithStack(err))
	}
}