package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"social-media/users/endpoint"

	transport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPServer(endpoints endpoint.Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(middleware)
	
	r.Methods("POST").Path("/register").Handler(transport.NewServer(
		endpoints.CreateUser,
		decodeCreateUserReq,
		encodeResponse,
	))

	r.Methods("POST").Path("/login").Handler(transport.NewServer(
		endpoints.Login,
		decodeLoginReq,
		encodeResponse,
	))

	r.Methods("GET").Path("/info/{login}").Handler(transport.NewServer(
		endpoints.GetUser,
		decodeGetUserReq,
		encodeResponse,
	))

	r.Methods("POST").Path("/info").Handler(transport.NewServer(
		endpoints.UpdateUser,
		decodeUpdateUserReq,
		encodeResponse,
	))

	r.Methods("POST").Path("/filter").Handler(transport.NewServer(
		endpoints.GetLoginsByInfo,
		decodeGetLoginsByInfoReq,
		encodeResponse,
	))

	r.Methods("POST").Path("/follow/{login}").Handler(transport.NewServer(
		endpoints.FollowUser,
		decodeFollowUserReq,
		encodeResponse,
	))

	r.Methods("GET").Path("/follow").Handler(transport.NewServer(
		endpoints.GetFollowedLogins,
		decodeGetFollowedLoginsReq,
		encodeResponse,
	))

	return r
}

func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
        next.ServeHTTP(w, r)
    })
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeCreateUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, err
	}

	return endpoint.CreateUserReq{
		Login:      r.FormValue("login"),
		FirstName:  r.FormValue("first_name"),
		SecondName: r.FormValue("second_name"),
		Password:   r.FormValue("passw1"),
	}, nil
}

func decodeLoginReq(_ context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseMultipartForm(10<<20); err != nil {
		return nil, err
	}

	return endpoint.LoginReq{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	}, nil
}

func decodeGetUserReq(_ context.Context, r *http.Request) (interface{}, error){
	params := mux.Vars(r)
	login, ok := params["login"]
	if !ok{
		return nil, errors.New("login not present")
	}
	return endpoint.GetUserReq{
		Login: login,
	}, nil
}

func decodeUpdateUserReq(_ context.Context, r *http.Request) (interface{}, error){
	token, err := r.Cookie("token")
	if err != nil{
		return nil, err
	}
	if err := r.ParseMultipartForm(10<<20); err != nil {
		return nil, err
	}

	return endpoint.UpdateUserReq{
		Token: token.Value,
		FirstName:  r.FormValue("first_name"),
		SecondName: r.FormValue("second_name"),
		Bio: r.FormValue("bio"),
		Interests: r.FormValue("interests"),
	}, nil
}

func decodeGetLoginsByInfoReq(_ context.Context, r *http.Request)(interface{}, error){
	if err := r.ParseMultipartForm(10<<20); err != nil {
		return nil, err
	}
	return endpoint.GetUserByInfoReq{
		Info: r.FormValue("info"),
	}, nil
}

func decodeFollowUserReq(_ context.Context, r *http.Request)(interface{}, error){
	token, err := r.Cookie("token")
	if err != nil{
		return nil, err
	}
	params := mux.Vars(r)
	login, ok := params["login"]
	if !ok{
		return nil, errors.New("login not present")
	}
	return endpoint.FollowUserReq{
		Token: token.Value,
		Login: login,
	}, nil
}

func decodeGetFollowedLoginsReq(_ context.Context, r *http.Request) (interface{}, error){
	token, err := r.Cookie("token")
	if err != nil{
		return nil, err
	}
	return endpoint.TokenReq{
		Token: token.Value,
	}, nil
}
