package handlers

import (
	"context"

	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
	pb "gitlab.mai.ru/cicada-chess/backend/user-service/pkg/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	userService interfaces.UserService
	pb.UnimplementedUserServiceServer
}

func NewGRPCHandler(userService interfaces.UserService) *GRPCHandler {
	return &GRPCHandler{
		userService: userService,
	}
}

func (h *GRPCHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	user, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserByEmailResponse{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Role:      int32(user.Role),
		Rating:    int32(user.Rating),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		IsActive:  user.IsActive,
	}, nil
}

func (h *GRPCHandler) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*pb.UpdateUserPasswordResponse, error) {
	err := h.userService.UpdatePasswordById(ctx, req.Id, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserPasswordResponse{Status: "success"}, nil
}
