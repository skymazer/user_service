package main

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/skymazer/user_service/broker"
	"github.com/skymazer/user_service/cache"
	"github.com/skymazer/user_service/loggerfx"
	"github.com/skymazer/user_service/middleware"
	pb "github.com/skymazer/user_service/proto"
	"github.com/skymazer/user_service/service"
	"github.com/skymazer/user_service/storage"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := loggerfx.New()
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	kafka, err := broker.New(logger, "user-service-logs")
	if err != nil {
		logger.Fatalf("failed to establish redis connection: %v", err)
	}
	defer kafka.Conn.Close()

	logger.SetStorager(kafka)

	redis, err := cache.New(logger)
	if err != nil {
		logger.Fatalf("failed to establish redis connection: %v", err)
	}
	defer redis.Conn.Close()

	dbUser, dbPassword, dbName :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")
	database, err := storage.New(dbUser, dbPassword, dbName)
	if err != nil {
		logger.Fatalf("failed to establish db connection: %v", err)
	}
	defer database.Conn.Close()

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(middleware.UserListCacheInterceptor(redis, logger),
			middleware.LoggerInterceptor(logger))),
	}
	grpcServer := grpc.NewServer(opts...)
	rpcServer, err := service.New(&database, logger)
	if err != nil {
		logger.Fatalf("failed to start rcp server: %v", err)
	}
	pb.RegisterUsersServer(grpcServer, rpcServer)
	go grpcServer.Serve(lis)

	logger.Info("Service started")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	logger.Info(fmt.Sprint(<-ch))
	logger.Info("Stopping API server.")
}
