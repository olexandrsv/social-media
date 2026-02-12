package server

import (
	"html/template"
	"net/http"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/log/service"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type httpServer struct {
	r       *mux.Router
	service service.Service
}

func newHttpServer(service service.Service) *httpServer {
	r := mux.NewRouter()
	srv := &httpServer{
		r:       r,
		service: service,
	}

	r.Methods("GET").Path("/logs").HandlerFunc(srv.getLogs)
	r.Methods("GET").Path("/health").HandlerFunc(common.NewHealthHandler("log"))

	return srv
}

func (s *httpServer) getLogs(w http.ResponseWriter, r *http.Request) {
	txt, err := s.service.GetLogs()
	if err != nil {
		common.WriteError(w, err)
		return
	}

	tmpl := template.Must(template.New("logs").Parse(`<pre> {{.}} </pre>`))
	tmpl.Execute(w, txt)
}

func (s *httpServer) run() {
	handler := common.CORS(s.r)
	if err := http.ListenAndServe(":"+config.App.LogService.HttpPort, handler); err != nil {
		log.Error(errors.WithStack(err))
	}
}
