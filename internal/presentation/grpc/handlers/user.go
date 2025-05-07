package handlers

import (
	"context"
	"errors"

	service "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	profileInterfaces "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/interfaces"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
	pb "gitlab.mai.ru/cicada-chess/backend/user-service/pkg/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	userService    interfaces.UserService
	profileService profileInterfaces.ProfileService
	pb.UnimplementedUserServiceServer
}

func NewGRPCHandler(userService interfaces.UserService, profileService profileInterfaces.ProfileService) *GRPCHandler {
	return &GRPCHandler{
		userService:    userService,
		profileService: profileService,
	}
}

func (h *GRPCHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	user, err := h.userService.Create(ctx, &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		IsActive: req.IsActive})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailExists) || errors.Is(err, service.ErrUsernameExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case errors.Is(err, service.ErrInvalidUUIDFormat) || errors.Is(err, entity.ErrPasswordTooShort) || errors.Is(err, entity.ErrInvalidEmail):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &pb.RegisterUserResponse{
		Id: user.ID,
	}, nil
}

func (h *GRPCHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	user, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
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
		switch err {
		case entity.ErrPasswordTooShort:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case service.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case service.ErrInvalidUUIDFormat:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.UpdateUserPasswordResponse{Status: "success"}, nil
}

func (h *GRPCHandler) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	user, err := h.userService.GetById(ctx, req.Id)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case service.ErrInvalidUUIDFormat:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.GetUserByIdResponse{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      int32(user.Role),
		Rating:    int32(user.Rating),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		IsActive:  user.IsActive,
	}, nil
}

func (h *GRPCHandler) ConfirmAccount(ctx context.Context, req *pb.ConfirmAccountRequest) (*pb.ConfirmAccountResponse, error) {
	err := h.userService.ConfirmAccount(ctx, req.Id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, service.ErrInvalidUUIDFormat):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	_, err = h.profileService.CreateProfile(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ConfirmAccountResponse{Status: "success"}, nil
}
