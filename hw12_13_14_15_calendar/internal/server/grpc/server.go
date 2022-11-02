package internalgrpc

import (
	"context"
	"net"
	"time"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
)

type Application interface {
	AddEvent(event *app.Event) error
	RemoveEvent(ID string) error
	UpdateEvent(event *app.Event) error
	GetEventsByDay(date time.Time) ([]app.Event, error)
	GetEventsByMonth(date time.Time) ([]app.Event, error)
	GetEventsByWeek(date time.Time) ([]app.Event, error)
}

type Server struct {
	server   *grpc.Server
	logger   Logger
	endpoint string
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

func NewServer(logger Logger, app Application, endpoint string) *Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingMiddleware(logger),
		),
	)
	service := &ServiceCalendar{logger: logger, app: app}
	pb.RegisterCalendarServer(server, service)
	return &Server{server, logger, endpoint}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("grpc server started on " + s.endpoint)
	lsn, err := net.Listen("tcp", s.endpoint)
	if err != nil {
		return err
	}
	if err := s.server.Serve(lsn); err != nil /*&& !errors.Is(err, http.ErrServerClosed)*/ {
		return err
	}
	s.logger.Info("grpc server stopped")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.server.Stop()
	return nil
}
