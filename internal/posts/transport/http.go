package transport

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/posts/endpoint"
	"strconv"

	transport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	e      endpoint.Endpoint
}

func newServer(e endpoint.Endpoint, r *mux.Router) *server {
	return &server{
		router: r,
		e:      e,
	}
}

func NewHTTPServer(e endpoint.Endpoint) *server {
	r := mux.NewRouter()
	s := newServer(e, r)

	r.Use(middleware)

	r.Methods("POST").Path("/users/posts").Handler(transport.NewServer(
		e.CreatePost,
		s.decodeCreatePostReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))
	r.Methods("GET").Path("/users/{id}/posts").Handler(transport.NewServer(
		e.GetPosts,
		s.decodeGetPostsReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))
	return s
}

func (s *server) Run() {
	err := http.ListenAndServe(":"+config.App.PostsService.Port, s.router)
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
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

func (s *server) decodeCreatePostReq(ctx context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(err)
		return nil, common.ErrNoToken
	}
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(err)
		return nil, common.ErrInvalidData
	}

	filesPath, err := processFormFiles(r.MultipartForm, "files[]")
	if err != nil {
		return nil, err
	}

	imagesPath, err := processFormFiles(r.MultipartForm, "images[]")
	if err != nil {
		return nil, err
	}

	return endpoint.CreatePostReq{
		Token:      token.Value,
		Text:       r.FormValue("text"),
		FilesPath:  filesPath,
		ImagesPath: imagesPath,
	}, nil
}

func processFileForm(form *multipart.Form, key string) ([]string, error) {
	filesHeaders := form.File[key]
	filesPaths := make([]string, 0, len(filesHeaders))
	for _, h := range filesHeaders {
		filesPaths = append(filesPaths, h.Filename)
		if err := saveFile(h); err != nil {
			return nil, err
		}
	}
	return filesPaths, nil
}

func saveFile(h *multipart.FileHeader) error {
	src, err := h.Open()
	if err != nil {
		return err
	}
	dst, err := os.OpenFile("./../../upload/"+h.Filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

func (s *server) decodeGetPostsReq(_ context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(err)
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
	return endpoint.GetPostsRequest{
		Token:  token.Value,
		UserID: id,
	}, nil
}

func (s *server) decodeTokenReq(_ context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(err)
		return nil, common.ErrNoToken
	}
	return endpoint.Token{
		Token: token.Value,
	}, nil
}

func (s *server) encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if response == nil {
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}
