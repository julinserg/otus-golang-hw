package internalgrpc

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func loggingMiddleware(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		h, err := handler(ctx, req)

		p, _ := peer.FromContext(ctx)
		ip := p.Addr.String()

		mD, _ := metadata.FromIncomingContext(ctx)
		ua := mD["user-agent"]

		rpcStatus := status.Convert(err)

		var sb strings.Builder
		sb.WriteString(ip + " ")
		sb.WriteString("[" + startTime.String() + "] ")
		sb.WriteString(info.FullMethod + " ")
		sb.WriteString(rpcStatus.Code().String() + " ")
		sb.WriteString(time.Since(startTime).String() + " ")
		sb.WriteString(`'` + strings.Join(ua, "") + `'`)
		logger.Info(sb.String())
		return h, err
	}
}
