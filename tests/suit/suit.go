package suit

import (
	"context"
	"net"
	"strconv"
	"testing"

	ssov1 "github.com/Touchme245/sso_protos/gen/go/sso"
	"github.com/Touchme245/sso_server/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcHost = "localhost"
)

type suit struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func NewSuit(t *testing.T) (context.Context, *suit) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPAth("../config/local.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Grpc.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	client, err := grpc.DialContext(context.Background(),
		grpcAdress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		t.Fatalf("grpc server connection failed: " + err.Error())
	}

	return ctx, &suit{
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(client),
	}

}

func grpcAdress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.Grpc.Port))
}
