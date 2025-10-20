package domain

import "github.com/hgyowan/go-email-grpc/domain/email"

type Repository interface {
	email.EmailRepository

	WithTransaction(fn func(txRepo Repository) error) error
}
