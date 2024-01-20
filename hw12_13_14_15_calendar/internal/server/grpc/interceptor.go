package internalgrpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

func UnaryInterceptor(logger Logger) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()

		clientIP := "-"
		userAgent := "-"

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if clientIPs := md.Get("x-forwarded-for"); len(clientIPs) > 0 {
				clientIP = clientIPs[0]
			}

			if userAgents := md.Get("User-Agent"); len(userAgents) > 0 {
				userAgent = userAgents[0]
			}
		}

		resp, err = handler(ctx, req)

		latency := time.Since(startTime)

		logger.Info(
			fmt.Sprintf(
				"%s [%s] %s %s \"%s\"",
				clientIP,
				time.Now().Format(time.RFC1123Z),
				info.FullMethod,
				latency,
				userAgent,
			),
		)

		return resp, err
	})
}
