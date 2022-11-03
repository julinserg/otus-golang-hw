package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var lastID uint64

func GetNextID() uint64 {
	curId := atomic.AddUint64(&lastID, 1)
	return curId
}

func GenerateActionId() string {
	curId := GetNextID()
	return fmt.Sprintf("%v:%v", time.Now().UTC().Format("20060102150405"), curId)
}

func main() {
	conn, err := grpc.Dial(":50001", grpc.WithInsecure(), grpc.WithUserAgent("user-agent"))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewCalendarClient(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		req, err := getRequest(reader)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}

		md := metadata.New(nil)
		md.Append("request_id", GenerateActionId())
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		if _, err := client.AddEvent(ctx, req); err != nil {
			log.Fatal(err)
		}

		log.Printf("event submitted")
	}
}

func getRequest(reader *bufio.Reader) (*pb.EventRequest, error) {
	log.Printf("write <event_id> <title>:")
	text, err := reader.ReadString('\n')
	if err != nil {
		return nil, errors.New("wrong input, try again")
	}

	parts := strings.Split(text, " ")
	if len(parts) < 2 {
		return nil, errors.New("wrong input, try again")
	}

	return &pb.EventRequest{
		Event: &pb.Event{
			Id:        parts[0],
			Title:     parts[1],
			TimeStart: timestamppb.Now(),
		},
	}, nil
}
