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

	r.Methods("PUT").Path("/users/posts/{id}").Handler(transport.NewServer(
		e.UpdatePost,
		s.decodeUpdatePostReq,
		s.encodeResponse,
		transport.ServerErrorEncoder(s.encodeError),
	))
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
	r.Methods("DELETE").Path("/users/posts/{id}").Handler(transport.NewServer(
		e.DeletePost,
		s.decodeDeletePostReq,
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

func processFormFiles(form *multipart.Form, key string) ([]string, error) {
	filesNames := []string{}
	files := form.File[key]
	for _, file := range files {
		filename, err := processFormFile(file)
		if err != nil {
			return nil, err
		}
		filesNames = append(filesNames, filename)
	}
	return filesNames, nil
}

func processFormFile(file *multipart.FileHeader) (string, error) {
	extension := filepath.Ext(file.Filename)
	id, err := uuid.NewUUID()
	if err != nil {
		log.Error(errors.WithStack(err))
		return "", common.ErrInternal
	}
	filename := id.String() + extension
	path := "./upload/" + filename
	if err := saveFile(file, path); err != nil {
		return "", err
	}
	return filename, nil
}

func saveFile(fileHeader *multipart.FileHeader, path string) error {
	file, err := fileHeader.Open()
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInvalidData
	}
	newFile, err := os.Create(path)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	_, err = io.Copy(newFile, file)
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
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
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}
	return endpoint.Token{
		Token: token.Value,
	}, nil
}

func (s *server) decodeUpdatePostReq(_ context.Context, r *http.Request) (interface{}, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.UpdatePostReq{
		Token:         token.Value,
		ID:            id,
		Text:          r.FormValue("text"),
		Images:        r.MultipartForm.File["images[]"],
		Files:         r.MultipartForm.File["files[]"],
		DeletedImages: r.MultipartForm.Value["deletedImages[]"],
		DeletedFiles:  r.MultipartForm.Value["deletedFiles[]"],
	}, nil
}

func (s *server) decodeDeletePostReq(_ context.Context, r *http.Request) (interface{}, error){
	token, err := r.Cookie("token")
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrNoToken
	}
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInvalidData
	}

	return endpoint.DeletePostReq{
		Token: token.Value,
		PostID: id,
	}, nil
}

func (s *server) encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if response == nil {
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}
