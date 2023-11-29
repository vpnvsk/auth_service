package auth

import (
	"auth_service/internal/services"
	authv1 "github.com/vpnvsk/protos_go/gen/auth"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
	services services.Services
}

func New(gRPC *grpc.Server, service services.Services) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{services: service})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is empty")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id is invalid")
	}
	token, err := s.services.LoginUser(ctx, req.GetLogin(), req.GetPassword(), int(req.GetAppId()))
	//TODO: handle exceptions more persistence
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authv1.LoginResponse{
		Token: token,
	}, err
}
func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is empty")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	uuid, err := s.services.RegisterUser(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authv1.RegisterResponse{Uuid: uuid}, nil
}
func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	if req.GetUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	isAdmin, err := s.services.UserIsAdmin(ctx, req.GetUuid(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authv1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
