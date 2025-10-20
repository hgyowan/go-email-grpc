package v1

import (
	"github.com/hgyowan/go-pkg-library/envs"
	pkgLibrary "github.com/hgyowan/go-pkg-library/grpc-library/grpc"
)

func EmailServiceClientProvider() EmailServiceClient {
	conn := pkgLibrary.MustNewGRPCClient(envs.MailGRPC)
	return NewEmailServiceClient(conn)
}
