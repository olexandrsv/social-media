package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/common/connection"
	"social-media/internal/users/endpoint"
	"social-media/internal/users/service"
	"strconv"
	"sync"

	"github.com/pkg/errors"

	transport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type httpServer struct {
	service   service.Service
	endpoints endpoint.Endpoints
	router    *mux.Router
}

func newHTTPServer(service service.Service, e endpoint.Endpoints, r *mux.Router) *httpServer {
	return &httpServer{
		service:   service,
		endpoints: e,
		router:    r,
	}
}

func NewHTTPServer(service service.Service, endpoints endpoint.Endpoints) *httpServer {
	r := mux.NewRouter()
	s := newHTTPServer(service, endpoints, r)

	r.Methods("POST").Path("/users").Handler(transport.NewServer(
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

	r.Methods("GET").Path("/users/followed").Handler(transport.NewServer(
		endpoints.GetFollowedUsers,
		s.decodeGetFollowedUsersReq,
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

	r.Methods("POST").Path("/users/{id}/followers").Handler(transport.NewServer(
		endpoints.FollowUser,
		s.decodeFollowUserReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))

	r.Methods("PUT").Path("/users/posts/{owner_id}/read").HandlerFunc(s.updateReadPosts)

	r.Path("/ws").HandlerFunc(s.upgradeToWS)

	return s
}

func (s *httpServer) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	handler := common.CORS(s.router)

	err := http.ListenAndServe(":"+config.App.Users.Service.HttpPort, handler)
	if err != nil {
		log.Error(err)
	}
}

func (s *httpServer) encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	code := 500
	msg := "Internal server error"
	if e, ok := err.(common.Error); ok {
		code = e.Code()
		msg = e.Message()
	}
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func (s *httpServer) encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if response == nil {
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func (s *httpServer) updateReadPosts(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		common.WriteError(w, common.ErrNoToken)
		return
	}

	params := mux.Vars(r)
	ownerIDParam, ok := params["owner_id"]
	if !ok {
		common.WriteError(w, common.ErrInvalidData)
		return
	}
	ownerID, err := strconv.Atoi(ownerIDParam)
	if err != nil {
		common.WriteError(w, common.ErrInvalidData)
		return
	}

	err = s.service.UpdateReadPosts(token.Value, ownerID, r.FormValue("last_read_post"))
	common.WriteResponse(w, nil, err)
}

func (s *httpServer) decodeCreateUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
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

func (s *httpServer) decodeLoginReq(_ context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.LoginReq{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	}, nil
}

func (s *httpServer) decodeGetUserReq(_ context.Context, r *http.Request) (interface{}, error) {
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

func (s *httpServer) decodeUpdateUserReq(_ context.Context, r *http.Request) (interface{}, error) {
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

func (s *httpServer) decodeGetUsersByInfoReq(_ context.Context, r *http.Request) (interface{}, error) {
	info := r.URL.Query().Get("info")

	return endpoint.GetUsersByInfoReq{
		Info: info,
	}, nil
}

func (s *httpServer) decodeFollowUserReq(_ context.Context, r *http.Request) (interface{}, error) {
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

	return endpoint.FollowUserReq{
		Token: token.Value,
		ID:    id,
	}, nil
}

func (s *httpServer) decodeGetFollowedUsersReq(_ context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}
	return endpoint.TokenReq{
		Token: token.Value,
	}, nil
}

func (s *httpServer) upgradeToWS(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrNoToken)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(errors.WithStack(err))
		writeError(w, common.ErrInternal)
		return
	}

	_ = conn
	_ = token

	s.service.SetUserConnection(token.Value, connection.NewGRPC(conn))
	//s.handleWebSocketConn(w, token.Value, conn)
}

func (s *httpServer) handleWebSocketConn(w http.ResponseWriter, token string, conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			log.Error(errors.WithStack(err))
			writeError(w, common.ErrInternal)
			return
		}

		log.Info("Received message: " + string(p))

		var t Notification
		if err := json.Unmarshal(p, &t); err != nil {
			log.Error(errors.WithStack(err))
			return
		}

		switch t.Type {
		case Message:
			var msg ReadMessagesNotification
			if err := json.Unmarshal(p, &msg); err != nil {
				log.Error(errors.WithStack(err))
				return
			}
			//s.service.UpdateReadMessages(token, msg.ChatID, msg.LastReadMessageID)
		case Post:
			var post ReadPostsNotification
			if err := json.Unmarshal(p, &post); err != nil {
				log.Error(errors.WithStack(err))
				return
			}
			//s.service.UpdateReadPosts(token, post.OwnerID, post.LastReadPostID)
		}

		//s.service.UpdateReadMessages(token, msg)
	}
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
