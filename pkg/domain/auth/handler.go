package auth

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mj23978/chat-backend/pkg/domain/user"
	m "github.com/mj23978/chat-backend/pkg/models"
	pb "github.com/mj23978/chat-backend/pkg/proto"
)

type authHandler struct {
	pb.UnimplementedAuthServiceServer
	userService user.UserService
	jwtManager  *JWTManager
}

func NewAuthHandler(userService user.UserService, jwtManager *JWTManager) *authHandler {
	return &authHandler{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

func (h *authHandler) Login(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var res *pb.LoginResponse
	switch req.GetReq().(type) {
	case *pb.LoginRequest_Email:
		user, err2 := h.userService.SignIn(req.GetEmail().Email, req.GetEmail().Password)
		if err2 != nil {
			return nil, status.Errorf(codes.NotFound, "Cannot Find Any User with Provided Information: %v", err2)
		}
		userRes, _ := turnUserToUserGrpc(user)
		token, err3 := h.jwtManager.Generate(userRes.Username)
		if err3 != nil {
			return nil, status.Error(codes.Internal, "Cannot Generate JWT Token")
		}
		refreshToken, err4 := h.jwtManager.GenerateRefresh(user.Username, user.Email, user.ID)
		if err4 != nil {
			return nil, status.Error(codes.Internal, "Cannot Generate Refresh Token")
		}
		res = &pb.LoginResponse{Token: token, RefreshToken: refreshToken}
	case *pb.LoginRequest_Token:
		res = &pb.LoginResponse{Token: "token", RefreshToken: "refreshToken"}
	}
	return res, nil
}

//func (h *authHandler) SignUp(_ context.Context, req *pb.SignUpRequest) (*pb.Token, error) {
//	if err := h.userService.SignUp(user); err != nil {
//		return nil, status.Errorf(codes.Unknown, "Something Happened: %v", err)
//	}
//	token, err2 := h.jwtManager.Generate(user.Username)
//	if err2 != nil {
//		return nil, status.Error(codes.Internal, "Cannot Generate JWT Token")
//	}
//	refreshToken, err3 := h.jwtManager.GenerateRefresh(user.Username, user.Email, user.ID)
//	if err3 != nil {
//		return nil, status.Error(codes.Internal, "Cannot Generate Refresh Token")
//	}
//	return &pb.Token{Token: token, RefreshToken: refreshToken}, nil
//}

func (h *authHandler) RefreshToken(_ context.Context, req *pb.RefreshTokenRequest) (*pb.Token, error) {
	var tokenRes, refreshTokenRes string
	switch req.GetId().(type) {
	case *pb.RefreshTokenRequest_RefreshToken:
		claims, err := h.jwtManager.VerifyRefresh(req.GetRefreshToken())
		if err != nil {
			log.Error().Msgf("%v", err)
			return nil, status.Error(codes.InvalidArgument, "Invalid Refresh Token")
		}
		token, err2 := h.jwtManager.Generate(claims.Username)
		if err2 != nil {
			return nil, status.Error(codes.Internal, "Cannot Generate JWT Token")
		}
		refreshToken, err3 := h.jwtManager.GenerateRefresh(claims.Username, claims.Email, claims.Id)
		if err3 != nil {
			return nil, status.Error(codes.Internal, "Cannot Generate Refresh Token")
		}
		tokenRes, refreshTokenRes = token, refreshToken
		break
	case *pb.RefreshTokenRequest_User:
		user, err := h.userService.SignIn(req.GetUser().Username, req.GetUser().Password)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Invalid Username or Password")
		}
		token, err2 := h.jwtManager.Generate(user.Username)
		if err2 != nil {
			return nil, status.Error(codes.Internal, "Cannot Generate JWT Token")
		}
		refreshToken, err3 := h.jwtManager.GenerateRefresh(user.Username, user.Email, user.ID)
		if err3 != nil {
			return nil, status.Error(codes.Internal, "Cannot Generate Refresh Token")
		}
		tokenRes, refreshTokenRes = token, refreshToken
		break
	}
	return &pb.Token{Token: tokenRes, RefreshToken: refreshTokenRes}, nil
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
