package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/users/endpoint"
	"strconv"

	"github.com/pkg/errors"

	transport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type server struct {
	endpoints endpoint.Endpoints
	router    *mux.Router
}

func newServer(e endpoint.Endpoints, r *mux.Router) *server {
	return &server{e, r}
}

func NewHTTPServer(endpoints endpoint.Endpoints) *server {
	r := mux.NewRouter()
	s := newServer(endpoints, r)

	r.Methods("POST").Path("/users").Queries().Handler(transport.NewServer(
		endpoints.CreateUser,
		s.decodeCreateUserReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	r.Methods("POST").Path("/users/login").Handler(transport.NewServer(
		endpoints.Login,
		s.decodeLoginReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	r.Methods("GET").Path("/users/{id}").Handler(transport.NewServer(
		endpoints.GetUser,
		s.decodeGetUserReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	r.Methods("PUT").Path("/users").Handler(transport.NewServer(
		endpoints.UpdateUser,
		s.decodeUpdateUserReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	r.Methods("GET").Path("/users").Queries("info", "{info}").Handler(transport.NewServer(
		endpoints.GetUsersByInfo,
		s.decodeGetUsersByInfoReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	r.Methods("POST").Path("/follow/{login}").Handler(transport.NewServer(
		endpoints.FollowUser,
		s.decodeFollowUserReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	r.Methods("GET").Path("/follow").Handler(transport.NewServer(
		endpoints.GetFollowedLogins,
		s.decodeGetFollowedLoginsReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	return s
}

func (s *server) Run() {
	handler := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"http://localhost:8080", "http://localhost:4200"}),
		handlers.AllowCredentials(),
	)(s.router)
	
	err := http.ListenAndServe(":"+config.App.UsersService.Port, handler)
	if err != nil {
		log.Error(err)
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

func (s *server) decodeCreateUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.CreateUserReq{
		Login:    r.FormValue("login"),
		Name:     r.FormValue("name"),
		Surname:  r.FormValue("surname"),
		Password: r.FormValue("password"),
	}, nil
}

func (s *server) decodeLoginReq(_ context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.LoginReq{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	}, nil
}

func (s *server) decodeGetUserReq(_ context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}

	params := mux.Vars(r)
	routeParam, ok := params["id"]
	if !ok {
		return nil, common.ErrInvalidData
	}

	id, err := strconv.Atoi(routeParam)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.GetUserReq{
		ID:    id,
		Token: token.Value,
	}, nil
}

func (s *server) decodeUpdateUserReq(_ context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.UpdateUserReq{
		Token:     token.Value,
		Name:      r.FormValue("name"),
		Surname:   r.FormValue("surname"),
		Bio:       r.FormValue("bio"),
		Interests: r.FormValue("interests"),
	}, nil
}

func (s *server) decodeGetUsersByInfoReq(_ context.Context, r *http.Request) (interface{}, error) {
	info := r.URL.Query().Get("info")

	return endpoint.GetUsersByInfoReq{
		Info: info,
	}, nil
}

func (s *server) decodeFollowUserReq(_ context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}
	params := mux.Vars(r)
	login, ok := params["login"]
	if !ok {
		return nil, common.ErrNoLogin
	}
	return endpoint.FollowUserReq{
		Token: token.Value,
		Login: login,
	}, nil
}

func (s *server) decodeGetFollowedLoginsReq(_ context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}
	return endpoint.TokenReq{
		Token: token.Value,
	}, nil
}
