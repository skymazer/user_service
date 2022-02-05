package rpc

import (
	"context"
	"errors"
	pb "github.com/skymazer/user_service/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
)

type Handler = pb.UsersServer

type user struct {
	Id   uint64
	Name string
	Mail string
}

type handler struct {
	pb.UnimplementedUsersServer

	users map[uint64]user
}

func New() (Handler, error) {
	var h handler
	h.users = make(map[uint64]user)
	return &h, nil
}

func (h *handler) AddUser(ctx context.Context, cur *pb.CreateUserReq) (*emptypb.Empty, error) {
	id := rand.Uint64()
	u := unmarshallUser(cur)
	u.Id = id
	h.users[id] = u
	return &emptypb.Empty{}, nil
}

func (h *handler) RemoveUser(ctx context.Context, u *pb.RemoveUserReq) (*emptypb.Empty, error) {

	if _, ok := h.users[u.Id]; !ok {
		return &emptypb.Empty{}, errors.New("User not found")
	}

	delete(h.users, u.Id)

	return &emptypb.Empty{}, nil
}

func (h *handler) ListUsers(context.Context, *emptypb.Empty) (*pb.ListUsersResp, error) {
	var resp pb.ListUsersResp
	resp.Users = make([]*pb.User, 0, len(h.users))

	for _, u := range h.users {
		resp.Users = append(resp.Users, marshallUser(&u))
	}

	return &resp, nil
}

func unmarshallUser(cur *pb.CreateUserReq) user {
	return user{
		Id:   0,
		Name: cur.Name,
		Mail: cur.Mail,
	}
}

func marshallUser(u *user) *pb.User {
	return &pb.User{
		Id:   u.Id,
		Name: u.Name,
		Mail: u.Mail,
	}
}
