package queue

import (
	"context"

	"github.com/hgyowan/go-email-grpc/domain"
)

type queueHandler struct {
	service               domain.Service
	externalQueueListener domain.ExternalQueueListener
}

func NewQueueHandler(service domain.Service, externalQueueListener domain.ExternalQueueListener) *queueHandler {
	h := &queueHandler{
		service:               service,
		externalQueueListener: externalQueueListener,
	}

	return h
}

func (q *queueHandler) Listen(ctx context.Context) {
	go listenEmailQueueHandler(ctx, q)
	<-ctx.Done()
}
