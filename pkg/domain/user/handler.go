package user

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mj23978/chat-backend-x/utils"
	m "github.com/mj23978/chat-backend/pkg/models"
	pb "github.com/mj23978/chat-backend/pkg/proto"
)

type userHandler struct {
	pb.UnimplementedUserServiceServer
	userService UserService
}

func NewUserHandler(userService UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) GetUser(_ context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	user, err := h.userService.FindById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Something Happened: %v", err)
	}
	userRes, _ := turnUserToUserGrpc(user)
	return userRes, nil
}

func (h *userHandler) UpdateUser(_ context.Context, req *pb.UpdateUserRequest) (*empty.Empty, error) {
	user, err := turnUserGrpcToUser(req.User)
	if err != nil {
		return nil, utils.ErrorWrap(err, "User.handler.UpdateUser")
	}
	if err := h.userService.UpdateUser(user); err != nil {
		return nil, status.Errorf(codes.Unknown, "Something Happened: %v", err)
	}
	return &empty.Empty{}, nil
}

func (h *userHandler) Status(_ context.Context, req *pb.StatusRequest) (*empty.Empty, error) {
	if err := h.userService.Status(req.Id, m.StatusType(req.Status)); err != nil {
		return nil, status.Errorf(codes.Unknown, "Something Happened: %v", err)
	}
	return &empty.Empty{}, nil
}

func turnUserToUserGrpc(user *m.User) (*pb.User, error) {
	userRes := &pb.User{}
	if err := copier.Copy(userRes, user); err != nil {
		return nil, err
	}
	userRes.Status = pb.StatusType(user.Status)
	userRes.Id = user.ID
	return userRes, nil
}

func turnUserGrpcToUser(user *pb.User) (*m.User, error) {
	userRes := &m.User{}
	if err := copier.Copy(userRes, user); err != nil {
		return nil, err
	}
	userRes.Status = m.StatusType(user.Status)
	userRes.ID = user.Id
	return userRes, nil
}
