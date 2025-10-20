package service

import (
	"context"
	"github.com/hgyowan/go-email-grpc/domain"
	"github.com/hgyowan/go-email-grpc/domain/email"
)

type service struct {
	email.EmailService

	repo domain.Repository

	externalRedisClient   domain.ExternalRedisClient
	externalQueueListener domain.ExternalQueueListener
	externalQueueEmitter  domain.ExternalQueueEmitter
	externalEmailSender   domain.ExternalEmailSender
	validator             domain.ExternalValidator
}

func NewService(ctx context.Context,
	repo domain.Repository,
	externalRedisClient domain.ExternalRedisClient,
	externalQueueListener domain.ExternalQueueListener,
	externalQueueEmitter domain.ExternalQueueEmitter,
	externalEmailSender domain.ExternalEmailSender,
	validator domain.ExternalValidator,
) domain.Service {
	s := &service{
		repo:                  repo,
		externalRedisClient:   externalRedisClient,
		externalQueueListener: externalQueueListener,
		externalQueueEmitter:  externalQueueEmitter,
		externalEmailSender:   externalEmailSender,
		validator:             validator,
	}
	s.register(ctx)
	return s
}

func (s *service) register(ctx context.Context) {
	registerEmailService(ctx, s)
}
