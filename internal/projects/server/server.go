package server

import (
	"encoding/json"
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"
	"social-media/internal/projects/domain/project"
	"social-media/internal/projects/service"
	"strconv"

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
		router:  r,
		service: srv,
	}
	r.Methods("GET").Path("/projects").HandlerFunc(s.getProjects)
	r.Methods("POST").Path("/projects").HandlerFunc(s.createProject)
	r.Methods("PUT").Path("/projects/{id}").HandlerFunc(s.updateProject)
	r.Methods("DELETE").Path("/projects/{id}").HandlerFunc(s.deleteProject)

	return s
}

func (s *server) Run() {
	handler := common.CORS(s.router)

	err := http.ListenAndServe(":"+config.App.Projects.Service.Port, handler)
	if err != nil {
		log.Error(err)
		panic(err)
	}

}

func (s *server) getProjects(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		writeError(w, common.ErrNoToken)
		return
	}

	projects, err := s.service.GetProjects(cookie.Value)
	if err != nil {
		writeError(w, err)
		return
	}

	models := slice.MustConvert[*project.Project, ProjectModel](projects, func(p *project.Project) ProjectModel {
		return ProjectModel{
			ID:          p.ID(),
			Login:       p.UserLogin(),
			Name:        p.Name(),
			Description: p.Description(),
			Stack:       p.Stack(),
		}
	})
	writeJSON(w, models)
}

func (s *server) createProject(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		writeError(w, common.ErrNoToken)
		return
	}

	if err := r.ParseMultipartForm(1 << 15); err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	stack := r.FormValue("stack")

	id, err := s.service.CreateProject(cookie.Value, name, description, stack)
	writeResponse(w, CreateProjectResp{id}, err)
}

func (s *server) updateProject(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		writeError(w, common.ErrNoToken)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		writeError(w, common.ErrInvalidData)
		return
	}

	projectID, err := strconv.Atoi(id)
	if err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	if err := r.ParseMultipartForm(1 << 15); err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	stack := r.FormValue("stack")

	err = s.service.UpdateProject(cookie.Value, projectID, name, description, stack)
	writeResponse(w, nil, err)
}

func (s *server) deleteProject(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		writeError(w, common.ErrNoToken)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		writeError(w, common.ErrInvalidData)
		return
	}

	projectID, err := strconv.Atoi(id)
	if err != nil {
		writeError(w, common.ErrInvalidData)
		return
	}

	err = s.service.DeleteProject(cookie.Value, projectID)
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
