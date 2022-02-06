package main

import (
	"fmt"
	"github.com/skymazer/user_service/cache"
	"github.com/skymazer/user_service/db"
	"github.com/skymazer/user_service/loggerfx"
	"github.com/skymazer/user_service/middleware"
	pb "github.com/skymazer/user_service/proto"
	rpcServer "github.com/skymazer/user_service/rpc"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	logger := loggerfx.ProvideLogger()
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	redis, err := cache.New(logger)
	if err != nil {
		logger.Fatalf("failed to establish redis connection: %v", err)
	}
	defer redis.Conn.Close()

	dbUser, dbPassword, dbName :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")
	database, err := db.New(dbUser, dbPassword, dbName)
	if err != nil {
		logger.Fatalf("failed to establish db connection: %v", err)
	}
	defer database.Conn.Close()

	opts := []grpc.ServerOption{
		middleware.UserListCacheInterceptor(redis, logger),
	}
	grpcServer := grpc.NewServer(opts...)
	rpcServer, err := rpcServer.New(&database)
	if err != nil {
		logger.Fatalf("failed to start rcp server: %v", err)
	}
	pb.RegisterUsersServer(grpcServer, rpcServer)
	go grpcServer.Serve(lis)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	logger.Info(fmt.Sprint(<-ch))
	logger.Info("Stopping API server.")
}
