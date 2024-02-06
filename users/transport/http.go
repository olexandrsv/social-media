package transport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"social-media/common"
	"social-media/users/endpoint"

	"github.com/pkg/errors"

	transport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type server struct {
	endpoints endpoint.Endpoints
	router *mux.Router
}

func newServer(e endpoint.Endpoints, r *mux.Router) *server {
	return &server{e, r}
}

func NewHTTPServer(endpoints endpoint.Endpoints) *server {
	r := mux.NewRouter()
	s := newServer(endpoints, r)

	r.Use(middleware)

	r.Methods("POST").Path("/register").Handler(transport.NewServer(
		endpoints.CreateUser,
		s.decodeCreateUserReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	r.Methods("POST").Path("/login").Handler(transport.NewServer(
		endpoints.Login,
		s.decodeLoginReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	r.Methods("GET").Path("/info/{login}").Handler(transport.NewServer(
		endpoints.GetUser,
		s.decodeGetUserReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	r.Methods("POST").Path("/info").Handler(transport.NewServer(
		endpoints.UpdateUser,
		s.decodeUpdateUserReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	r.Methods("POST").Path("/filter").Handler(transport.NewServer(
		endpoints.GetLoginsByInfo,
		s.decodeGetLoginsByInfoReq,
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
	err := http.ListenAndServe(":8081", s.router)
	if err != nil {
		log.Fatal(err)
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}

func (s *server) encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	s.Log(err)
	code := 500
	msg := "Internal server error"
	if e, ok := err.(common.Error); ok {
		code = e.Code
		msg = e.Message
	}
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func (s *server) encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func (s *server) decodeCreateUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		s.Log(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.CreateUserReq{
		Login:      r.FormValue("login"),
		FirstName:  r.FormValue("first_name"),
		SecondName: r.FormValue("second_name"),
		Password:   r.FormValue("passw1"),
	}, nil
}

func (s *server) decodeLoginReq(_ context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		s.Log(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.LoginReq{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	}, nil
}

func (s *server) decodeGetUserReq(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	login, ok := params["login"]
	if !ok {
		return nil, common.ErrNoLogin
	}
	return endpoint.GetUserReq{
		Login: login,
	}, nil
}

func (s *server) decodeUpdateUserReq(_ context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		s.Log(errors.WithStack(err))
		return nil, common.ErrNoToken
	}
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		s.Log(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.UpdateUserReq{
		Token:      token.Value,
		FirstName:  r.FormValue("first_name"),
		SecondName: r.FormValue("second_name"),
		Bio:        r.FormValue("bio"),
		Interests:  r.FormValue("interests"),
	}, nil
}

func (s *server) decodeGetLoginsByInfoReq(_ context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		s.Log(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	return endpoint.GetUserByInfoReq{
		Info: r.FormValue("info"),
	}, nil
}

func (s *server) decodeFollowUserReq(_ context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		s.Log(errors.WithStack(err))
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
		s.Log(errors.WithStack(err))
		return nil, common.ErrNoToken
	}
	return endpoint.TokenReq{
		Token: token.Value,
	}, nil
}

func (s *server) Log(err error){
	s.endpoints.Log(err)
}
