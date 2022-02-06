package middleware

import (
	"context"
	"github.com/skymazer/user_service/loggerfx"
	"google.golang.org/grpc"
)

func LoggerInterceptor(l *loggerfx.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		l.Info("API call to ", info.FullMethod)

		h, err := handler(ctx, req)

		return h, err
	}
}
