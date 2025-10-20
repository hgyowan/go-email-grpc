package external

import (
	"github.com/hgyowan/go-email-grpc/domain"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgGrpc "github.com/hgyowan/go-pkg-library/grpc-library/grpc"
)

type externalGrpcServer struct {
	server pkgGrpc.GrpcServer
	port   string
}

func (g *externalGrpcServer) Server() pkgGrpc.GrpcServer {
	return g.server
}

func (g *externalGrpcServer) Port() string {
	return g.port
}

func MustNewGRPCServer() domain.ExternalGRPCServer {
	server := pkgGrpc.MustNewGRPCServer()

	return &externalGrpcServer{
		server: server,
		port:   envs.ServerPort,
	}
}
