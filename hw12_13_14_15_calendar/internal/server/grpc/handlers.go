package internalgrpc

import (
	"context"
	"time"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServiceCalendar struct {
	pb.UnimplementedCalendarServer
	logger Logger
	app    Application
}

func (s *ServiceCalendar) AddEvent(ctx context.Context, req *pb.EventRequest) (*pb.ErrorResponse, error) {
	err := s.app.AddEvent(&app.EventApp{ID: req.Event.Id, Title: req.Event.Title, Description: req.Event.Description,
		UserID: req.Event.UserID, NotificationTime: time.Duration(req.Event.NotificationTime), TimeStart: req.Event.TimeStart.AsTime(),
		TimeEnd: req.Event.TimeEnd.AsTime()})
	return &pb.ErrorResponse{}, err
}

func (s *ServiceCalendar) RemoveEvent(ctx context.Context, req *pb.IdRequest) (*pb.ErrorResponse, error) {
	err := s.app.RemoveEvent(req.Id)
	return &pb.ErrorResponse{}, err
}

func (s *ServiceCalendar) UpdateEvent(ctx context.Context, req *pb.EventRequest) (*pb.ErrorResponse, error) {
	err := s.app.UpdateEvent(&app.EventApp{ID: req.Event.Id, Title: req.Event.Title, Description: req.Event.Description,
		UserID: req.Event.UserID, NotificationTime: time.Duration(req.Event.NotificationTime), TimeStart: req.Event.TimeStart.AsTime(),
		TimeEnd: req.Event.TimeEnd.AsTime()})
	return &pb.ErrorResponse{}, err
}

func (s *ServiceCalendar) GetEventsByDay(ctx context.Context, req *pb.TimeRequest) (*pb.EventsResponse, error) {
	events, err := s.app.GetEventsByDay(req.Date.AsTime())
	resp := &pb.EventsResponse{}
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	resp.Events = make([]*pb.Event, 0, len(events))
	for _, e := range events {
		ev := &pb.Event{Id: e.ID, Title: e.Title, Description: e.Description,
			UserID: e.UserID, NotificationTime: int64(e.NotificationTime), TimeStart: timestamppb.New(e.TimeStart),
			TimeEnd: timestamppb.New(e.TimeEnd)}
		resp.Events = append(resp.Events, ev)
	}
	return resp, err

}
func (s *ServiceCalendar) GetEventsByMonth(ctx context.Context, req *pb.TimeRequest) (*pb.EventsResponse, error) {
	events, err := s.app.GetEventsByMonth(req.Date.AsTime())
	resp := &pb.EventsResponse{}
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}

	resp.Events = make([]*pb.Event, 0, len(events))
	for _, e := range events {
		resp.Events = append(resp.Events, &pb.Event{Id: e.ID, Title: e.Title, Description: e.Description,
			UserID: e.UserID, NotificationTime: int64(e.NotificationTime), TimeStart: timestamppb.New(e.TimeStart),
			TimeEnd: timestamppb.New(e.TimeEnd)})
	}
	return resp, err
}
func (s *ServiceCalendar) GetEventsByWeek(ctx context.Context, req *pb.TimeRequest) (*pb.EventsResponse, error) {
	events, err := s.app.GetEventsByWeek(req.Date.AsTime())
	resp := &pb.EventsResponse{}
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}

	resp.Events = make([]*pb.Event, 0, len(events))
	for _, e := range events {
		resp.Events = append(resp.Events, &pb.Event{Id: e.ID, Title: e.Title, Description: e.Description,
			UserID: e.UserID, NotificationTime: int64(e.NotificationTime), TimeStart: timestamppb.New(e.TimeStart),
			TimeEnd: timestamppb.New(e.TimeEnd)})
	}
	return resp, err
}
