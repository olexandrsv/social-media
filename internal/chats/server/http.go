package server

import (
	"encoding/json"
	"net/http"
	"social-media/internal/chats/service"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server interface {
	Run()
}

type httpServer struct {
	router  *mux.Router
	service service.Service
}

func newHttpServer(srv service.Service) Server {
	r := mux.NewRouter()
	s := &httpServer{
		router:  r,
		service: srv,
	}
	r.Methods("GET").Path("/chats").HandlerFunc(s.getChats)
	r.Methods("POST").Path("/chats").HandlerFunc(s.createChat)
	r.Methods("GET").Path("/chats/{id}").HandlerFunc(s.getChat)
	r.Methods("PUT").Path("/chats/{id}").HandlerFunc(s.updateChat)
	r.Methods("DELETE").Path("/chats/{id}").HandlerFunc(s.deleteChat)
	r.Methods("PUT").Path("/users/chats/{chat_id}/read").HandlerFunc(s.updateReadMessages)

	return s
}

func (s *httpServer) Run() {
	handler := common.CORS(s.router)

	err := http.ListenAndServe(":"+config.App.Chats.Service.HttpPort, handler)
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

func (s *httpServer) updateReadMessages(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		common.WriteError(w, common.ErrNoToken)
		return
	}

	params := mux.Vars(r)
	chatIDParam, ok := params["chat_id"]
	if !ok {
		common.WriteError(w, common.ErrInvalidData)
		return
	}
	chatID, err := strconv.Atoi(chatIDParam)
	if err != nil {
		common.WriteError(w, common.ErrInvalidData)
		return
	}

	err = s.service.UpdateReadMessages(token.Value, chatID, r.FormValue("last_read_message"))
	common.WriteResponse(w, nil, err)
}

func (s *httpServer) getChats(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		writeError(w, common.ErrNoToken)
		return
	}

	token := cookie.Value
	chats, err := s.service.UserChats(token)
	if err != nil {
		writeError(w, common.ErrInternal)
		return
	}

	models := make([]BasicUserChat, 0, len(chats))
	for _, chat := range chats {
		model := BasicUserChat{
			ID:   chat.ID(),
			Name: chat.Name(),
		}
		models = append(models, model)
	}

	writeJSON(w, UserChatsResp(models))
}

func (s *httpServer) createChat(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		writeError(w, common.ErrNoToken)
		return
	}

	if err := r.ParseMultipartForm(1 << 15); err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	ids, err := slice.Convert(r.MultipartForm.Value["users_ids[]"], func(rawID string) (int, error) {
		return strconv.Atoi(rawID)
	})
	if err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	chat, err := s.service.CreateChat(service.CreateChatReq{
		Token:    cookie.Value,
		Name:     r.FormValue("name"),
		UsersIDs: ids,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, chatToModel(chat))
}

func (s *httpServer) getChat(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		writeError(w, common.ErrNoToken)
		return
	}

	params := mux.Vars(r)
	rawID, ok := params["id"]
	if !ok {
		writeError(w, common.ErrInvalidData)
		return
	}

	id, err := strconv.Atoi(rawID)
	if err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	chat, err := s.service.Chat(token.Value, id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, chatToModel(chat))
}

func (s *httpServer) updateChat(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		writeError(w, common.ErrNoToken)
		return
	}

	if err := r.ParseMultipartForm(1 << 15); err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	ids, err := slice.Convert(r.MultipartForm.Value["users_ids[]"], func(rawID string) (int, error) {
		return strconv.Atoi(rawID)
	})
	if err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	rawID, ok := mux.Vars(r)["id"]
	if !ok {
		writeError(w, common.ErrInvalidData)
		return
	}
	id, err := strconv.Atoi(rawID)
	if err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	err = s.service.UpdateChat(service.UpdateChatReq{
		Token:    token.Value,
		ID:       id,
		Name:     r.FormValue("name"),
		UsersIDs: ids,
	})
	writeResponse(w, nil, err)
}

func (s *httpServer) deleteChat(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		writeError(w, common.ErrNoToken)
		return
	}

	rawID, ok := mux.Vars(r)["id"]
	if !ok {
		writeError(w, common.ErrInvalidData)
		return
	}
	id, err := strconv.Atoi(rawID)
	if err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	err = s.service.DeleteChat(token.Value, id)
	writeResponse(w, nil, err)
}

func writeResponse(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, resp)
}

func writeError(w http.ResponseWriter, err error) {
	code := 500
	msg := "Internal server error"
	if e, ok := err.(common.Error); ok {
		code = e.Code()
		msg = e.Message()
	}
	w.WriteHeader(code)

	_, err = w.Write([]byte(msg))
	if err != nil {
		log.Error(errors.WithStack(err))
	}
}

func writeJSON(w http.ResponseWriter, resp interface{}) {
	w.WriteHeader(200)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error(errors.WithStack(err))
	}
}
