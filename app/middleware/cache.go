package middleware

import (
	"context"
	errors2 "errors"
	"github.com/pkg/errors"
	"github.com/skymazer/user_service/loggerfx"
	pb "github.com/skymazer/user_service/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"strings"
)

type Cacher interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	Remove(key string) error
}

var ErrNotFound = errors.New("no matching record")

const cachedService = "/users.Users/"
const cachedMethod = cachedService + "ListUsers"

func UserListCacheInterceptor(r Cacher, l *loggerfx.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		if info.FullMethod == cachedMethod {
			data, err := r.Get(cachedMethod)
			if err == nil {
				var s pb.ListUsersResp
				if err = proto.Unmarshal(data, &s); err != nil {
					l.Warnf("failed to unmarshall cache: %v", err)
				}
				return &s, nil
			}
			if !errors2.Is(err, ErrNotFound) {
				l.Warnf("failed to fetch cache: %v", err)
			}
		}

		h, err := handler(ctx, req)
		if err != nil {
			return h, err
		}

		if info.FullMethod == cachedMethod {
			bytes, err := proto.Marshal(h.(proto.Message))
			if err != nil {
				l.Warnf("failed to marshall responce: %v", err)
				return h, err
			}
			if err = r.Set(cachedMethod, bytes); err != nil {
				l.Warnf("failed to cache response: %v", err)
			}
		} else if strings.HasPrefix(info.FullMethod, cachedService) {
			if err := r.Remove(cachedMethod); err != nil {
				if !errors2.Is(err, ErrNotFound) {
					l.Warnf("failed to invalidate cache: %v", err)
				}
			}
		}

		return h, err
	}
}
