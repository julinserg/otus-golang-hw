package internalhttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app"
)

type Application interface {
	AddEvent(event *app.EventApp) error
	RemoveEvent(ID string) error
}

type Server struct { // TODO
	server   *http.Server
	logger   Logger
	endpoint string
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func NewServer(logger Logger, app Application, endpoint string) *Server {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    endpoint,
		Handler: loggingMiddleware(mux, logger),
	}
	ch := calendarHandler{logger, app}
	mux.HandleFunc("/", hellowHandler)
	mux.HandleFunc("/add", ch.addEvent)
	mux.HandleFunc("/remove", ch.removeEvent)
	return &Server{server, logger, endpoint}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("server started on " + s.endpoint)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	s.logger.Info("server stopped")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
