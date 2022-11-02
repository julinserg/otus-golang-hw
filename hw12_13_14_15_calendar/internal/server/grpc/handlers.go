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

func eventFromPb(ev *pb.Event) *app.Event {
	return &app.Event{
		ID: ev.Id, Title: ev.Title, Description: ev.Description,
		UserID: ev.UserID, NotificationTime: time.Duration(ev.NotificationTime), TimeStart: ev.TimeStart.AsTime(),
		TimeEnd: ev.TimeEnd.AsTime(),
	}
}

func eventToPb(ev *app.Event) *pb.Event {
	return &pb.Event{
		Id: ev.ID, Title: ev.Title, Description: ev.Description,
		UserID: ev.UserID, NotificationTime: int64(ev.NotificationTime), TimeStart: timestamppb.New(ev.TimeStart),
		TimeEnd: timestamppb.New(ev.TimeEnd),
	}
}

func fillEventsResponse(events []app.Event, err error) (*pb.EventsResponse, error) {
	resp := &pb.EventsResponse{}
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	resp.Events = make([]*pb.Event, 0, len(events))
	for _, e := range events {
		resp.Events = append(resp.Events, eventToPb(&e))
	}
	return resp, nil
}

func fillErrorResponse(err error) (*pb.ErrorResponse, error) {
	resp := &pb.ErrorResponse{}
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	return resp, nil
}

func (s *ServiceCalendar) AddEvent(ctx context.Context, req *pb.EventRequest) (*pb.ErrorResponse, error) {
	err := s.app.AddEvent(eventFromPb(req.Event))
	return fillErrorResponse(err)
}

func (s *ServiceCalendar) RemoveEvent(ctx context.Context, req *pb.IdRequest) (*pb.ErrorResponse, error) {
	err := s.app.RemoveEvent(req.Id)
	return fillErrorResponse(err)
}

func (s *ServiceCalendar) UpdateEvent(ctx context.Context, req *pb.EventRequest) (*pb.ErrorResponse, error) {
	err := s.app.UpdateEvent(eventFromPb(req.Event))
	return fillErrorResponse(err)
}

func (s *ServiceCalendar) GetEventsByDay(ctx context.Context, req *pb.TimeRequest) (*pb.EventsResponse, error) {
	events, err := s.app.GetEventsByDay(req.Date.AsTime())
	return fillEventsResponse(events, err)
}

func (s *ServiceCalendar) GetEventsByMonth(ctx context.Context, req *pb.TimeRequest) (*pb.EventsResponse, error) {
	events, err := s.app.GetEventsByMonth(req.Date.AsTime())
	return fillEventsResponse(events, err)
}

func (s *ServiceCalendar) GetEventsByWeek(ctx context.Context, req *pb.TimeRequest) (*pb.EventsResponse, error) {
	events, err := s.app.GetEventsByWeek(req.Date.AsTime())
	return fillEventsResponse(events, err)
}
