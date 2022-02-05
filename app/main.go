package main

import (
	"context"
	"github.com/skymazer/user_service/loggerfx"
	pb "github.com/skymazer/user_service/proto"
	rpcServer "github.com/skymazer/user_service/rpc"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	fx.New(
		fx.Provide(rpcServer.New),
		fx.Invoke(registerHooks),
		loggerfx.Module,
	).Run()
}

func registerHooks(
	lifecycle fx.Lifecycle, logger *zap.SugaredLogger, rpcServer rpcServer.Handler,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {

				lis, err := net.Listen("tcp", ":8081")
				if err != nil {
					logger.Fatalf("failed to listen: %v", err)
				}
				var opts []grpc.ServerOption
				grpcServer := grpc.NewServer(opts...)
				pb.RegisterUsersServer(grpcServer, rpcServer)
				go grpcServer.Serve(lis)

				return nil
			},
			OnStop: func(context.Context) error {
				return logger.Sync()
			},
		},
	)
}
