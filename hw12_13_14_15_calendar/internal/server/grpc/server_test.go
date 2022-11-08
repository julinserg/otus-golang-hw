package internalgrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app_calendar"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/server/grpc/pb"
	memorystorage "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LoggerFakeImpl struct{}

func (l LoggerFakeImpl) Info(msg string) {
	fmt.Println("Info: " + msg)
}

func (l LoggerFakeImpl) Error(msg string) {
	fmt.Println("Error: " + msg)
}

func (l LoggerFakeImpl) Debug(msg string) {
	fmt.Println("Debug: " + msg)
}

func (l LoggerFakeImpl) Warn(msg string) {
	fmt.Println("Warn: " + msg)
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestServiceGRPC(t *testing.T) {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	logg := LoggerFakeImpl{}
	storage := memorystorage.New()
	calendar := app_calendar.New(logg, storage)
	service := &ServiceCalendar{logger: logg, app: calendar}
	pb.RegisterCalendarServer(s, service)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalendarClient(conn)
	_, err = client.AddEvent(ctx, &pb.EventRequest{Event: &pb.Event{
		Id:        "1",
		Title:     "Title1",
		TimeStart: timestamppb.New(time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC)),
	}})
	require.Nil(t, err)

	_, err = client.AddEvent(ctx, &pb.EventRequest{Event: &pb.Event{
		Id:        "2",
		Title:     "Title1",
		TimeStart: timestamppb.New(time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC)),
	}})
	require.NotNil(t, err)
	require.Equal(t, err.Error(), "rpc error: code = Unknown desc = time event is busy")

	_, err = client.AddEvent(ctx, &pb.EventRequest{Event: &pb.Event{
		Id:        "2",
		Title:     "Title2",
		TimeStart: timestamppb.New(time.Date(2022, time.Month(1), 3, 1, 10, 30, 0, time.UTC)),
	}})
	require.Nil(t, err)

	_, err = client.AddEvent(ctx, &pb.EventRequest{Event: &pb.Event{
		Id:        "3",
		Title:     "Title3",
		TimeStart: timestamppb.New(time.Date(2022, time.Month(1), 14, 1, 10, 30, 0, time.UTC)),
	}})
	require.Nil(t, err)

	resp, err := client.GetEventsByDay(ctx, &pb.TimeRequest{Date: timestamppb.New(time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC))})
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, resp.Events[0].Id, "1")
	require.Equal(t, resp.Events[0].Title, "Title1")

	resp, err = client.GetEventsByWeek(ctx, &pb.TimeRequest{Date: timestamppb.New(time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC))})
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, len(resp.Events), 2)

	resp, err = client.GetEventsByMonth(ctx, &pb.TimeRequest{Date: timestamppb.New(time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC))})
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, len(resp.Events), 3)

	_, err = client.RemoveEvent(ctx, &pb.IdRequest{Id: "3"})
	require.Nil(t, err)

	resp, err = client.GetEventsByMonth(ctx, &pb.TimeRequest{Date: timestamppb.New(time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC))})
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, len(resp.Events), 2)

	_, err = client.UpdateEvent(ctx, &pb.EventRequest{Event: &pb.Event{
		Id:        "1",
		Title:     "Title1",
		TimeStart: timestamppb.New(time.Date(2025, time.Month(1), 1, 1, 10, 30, 0, time.UTC)),
	}})
	require.Nil(t, err)

	resp, err = client.GetEventsByMonth(ctx, &pb.TimeRequest{Date: timestamppb.New(time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC))})
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, len(resp.Events), 1)
	require.Equal(t, resp.Events[0].Id, "2")
	require.Equal(t, resp.Events[0].Title, "Title2")
}
