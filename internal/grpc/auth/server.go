package auth

import (
	"context"

	ssov1 "github.com/Touchme245/sso_protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverApi struct {
	ssov1.UnimplementedAuthServer
}

func Register(grpServer *grpc.Server) {
	ssov1.RegisterAuthServer(grpServer, &serverApi{})
}

func (s *serverApi) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	panic("not implemented")
}

func (s *serverApi) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("not implemented")
}

func (s *serverApi) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("not implemented")
}
