package service

import (
	"context"
	"encoding/json"
	"github.com/skymazer/user_service/loggerfx"
	m "github.com/skymazer/user_service/models"
	pb "github.com/skymazer/user_service/proto"
	"github.com/skymazer/user_service/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type Handler = pb.UsersServer

type Storage interface {
	AddUser(item *m.User) (m.IdType, error)
	DeleteUser(userId m.IdType) error
	GetAllUsers() ([]*m.User, error)
}

type handler struct {
	pb.UnimplementedUsersServer
	database *storage.Database
	log      *loggerfx.Logger
}

func New(database *storage.Database, log *loggerfx.Logger) (Handler, error) {
	var h handler
	h.database = database
	h.log = log
	return &h, nil
}

func (h *handler) AddUser(ctx context.Context, cur *pb.CreateUserReq) (*emptypb.Empty, error) {
	res := &emptypb.Empty{}
	u := unmarshallUser(cur)

	if len(u.Name) == 0 {
		return res, status.Error(codes.InvalidArgument, "Mail must not be null")
	}
	if len(u.Mail) == 0 {
		return res, status.Error(codes.InvalidArgument, "Name must not be null")
	}

	id, err := h.database.AddUser(&u)
	if err != nil {
		return res, status.Errorf(codes.Internal,
			"Failed to create user: %v", err)
	}

	logEntry, err := json.Marshal(struct {
		Timestamp uint64 `json:"timestamp"`
		UserId    uint64 `json:"userId"`
		Name      string `json:"name"`
	}{uint64(time.Now().Unix()), uint64(id), u.Name})
	h.log.InfoS(logEntry)

	return &emptypb.Empty{}, nil
}

func (h *handler) RemoveUser(ctx context.Context, u *pb.RemoveUserReq) (*emptypb.Empty, error) {

	if err := h.database.DeleteUser(m.IdType(u.Id)); err != nil {
		switch err {
		case storage.ErrNoMatch:
			return &emptypb.Empty{}, status.Error(codes.NotFound,
				"User doesn't exist")
		default:
			return &emptypb.Empty{}, status.Errorf(codes.Internal,
				"Failed to remove user: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (h *handler) ListUsers(context.Context, *emptypb.Empty) (*pb.ListUsersResp, error) {
	var resp pb.ListUsersResp

	stored, err := h.database.GetAllUsers()
	if err != nil {
		return &resp, status.Errorf(codes.Internal,
			"Failed to fetch user: %v", err)
	}

	for _, u := range stored {
		resp.Users = append(resp.Users, marshallUser(u))
	}

	return &resp, nil
}

func unmarshallUser(cur *pb.CreateUserReq) m.User {
	return m.User{
		Id:   0,
		Name: cur.Name,
		Mail: cur.Mail,
	}
}

func marshallUser(u *m.User) *pb.User {
	return &pb.User{
		Id:   uint64(u.Id),
		Name: u.Name,
		Mail: u.Mail,
	}
}
