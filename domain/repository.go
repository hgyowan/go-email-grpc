package domain

import (
	"context"
	"github.com/hgyowan/go-email-grpc/domain/email"
)

type Repository interface {
	email.EmailRepository

	WithTransaction(ctx context.Context, fn func(txRepo Repository) error) error
}
