package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/posts/endpoint"
	"social-media/internal/posts/service"

	transport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type server struct {
	router  *mux.Router
	e       endpoint.Endpoint
	service service.Service
}

func newServer(e endpoint.Endpoint, service service.Service, r *mux.Router) *server {
	return &server{
		router:  r,
		e:       e,
		service: service,
	}
}

func NewHTTPServer(e endpoint.Endpoint, service service.Service) *server {
	r := mux.NewRouter()
	s := newServer(e, service, r)

	r.Methods("PUT").Path("/users/posts/{id}").HandlerFunc(s.updatePost)
	r.Methods("POST").Path("/users/posts").HandlerFunc(s.createPost)
	r.Methods("GET").Path("/users/{id}/posts").HandlerFunc(s.getPosts)
	r.Methods("DELETE").Path("/users/posts/{id}").HandlerFunc(s.deletePost)

	r.Methods("GET").Path("/posts/{id}/comments").HandlerFunc(s.postComments)

	r.Methods("POST").Path("/users/posts/{id}/comments").Handler(transport.NewServer(
		e.CreatePostComment,
		s.decodeCreatePostCommentReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))
	r.Methods("DELETE").Path("/users/posts/{post_id}/comments/{comment_id}").Handler(transport.NewServer(
		e.DeletePostComment,
		s.decodeDeletePostCommentReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))
	r.Methods("PUT").Path("/users/posts/comments/{id}").Handler(transport.NewServer(
		e.UpdateComment,
		s.decodeUpdateCommentReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))
	r.Methods("GET").Path("/posts/comments/{id}/comments").Handler(transport.NewServer(
		e.CommentComments,
		s.decodeCommentCommentsReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))
	r.Methods("POST").Path("/users/posts/comments/{id}/comments").Handler(transport.NewServer(
		e.CreateCommentComment,
		s.decodeCreateCommnetCommentReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))
	r.Methods("DELETE").Path("/users/posts/comments/{parent_id}/comments/{comment_id}").Handler(transport.NewServer(
		e.DeleteCommentComment,
		s.decodeDeleteCommentCommentReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	r.Methods("GET").Path("/chats/{chat_id}/messages").HandlerFunc(s.getChatMessages)
	r.Methods("POST").Path("/chats/{chat_id}/messages").HandlerFunc(s.createChatMessage)
	r.Methods("PUT").Path("/chats/messages/{id}").HandlerFunc(s.updateChatMessage)
	r.Methods("DELETE").Path("/chats/messages/{id}").HandlerFunc(s.deleteChatMessage)

	return s
}

func (s *server) Run() {
	handler := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"http://localhost:8080", "http://localhost:4200"}),
		handlers.AllowCredentials(),
	)(s.router)

	err := http.ListenAndServe(":"+config.App.PostsService.Port, handler)
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

func (s *server) encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	code := 500
	msg := "Internal server error"
	if e, ok := err.(common.Error); ok {
		code = e.Code()
		msg = e.Message()
	}
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func (s *server) encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if response == nil {
		return nil
	}
	return json.NewEncoder(w).Encode(response)
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
