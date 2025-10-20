package grpc

import (
	"context"
	"github.com/hgyowan/go-email-grpc/domain"
	emailV1 "github.com/hgyowan/go-email-grpc/gen/email/v1"
)

type grpcHandler struct {
	emailV1.EmailServiceServer

	service            domain.Service
	externalGRPCServer domain.ExternalGRPCServer
}

func NewGRPCHandler(service domain.Service, externalGRPCServer domain.ExternalGRPCServer) *grpcHandler {
	h := &grpcHandler{
		service:            service,
		externalGRPCServer: externalGRPCServer,
	}

	h.register()

	return h
}

func (h *grpcHandler) Listen(ctx context.Context) {
	h.externalGRPCServer.Server().Serve(ctx, h.externalGRPCServer.Port())
}

func (h *grpcHandler) register() {
	registerEmailGRPCHandler(h)
}
