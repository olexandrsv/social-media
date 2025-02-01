package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-media/internal/chats/service"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server interface {
	Run()
}

type server struct {
	router  *mux.Router
	service service.Service
}

func New(srv service.Service) Server {
	r := mux.NewRouter()
	s := &server{
		router: r,
		service: srv,
	}
	r.Methods("GET").Path("/chats").HandlerFunc(s.getChats)

	return s
}

func (s *server) getChats(w http.ResponseWriter, r *http.Request) {
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

	log.Error(fmt.Errorf("%#v", models))
	writeJSON(w, UserChatsResp(models))
}

func (s *server) Run() {
	handler := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"http://localhost:8080", "http://localhost:4200"}),
		handlers.AllowCredentials(),
	)(s.router)

	err := http.ListenAndServe(":"+config.App.ChatsService.Port, handler)
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

func writeResponse(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		writeError(w, err)
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
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error(errors.WithStack(err))
	}
}
