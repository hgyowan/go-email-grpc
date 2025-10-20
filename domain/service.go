package domain

import "github.com/hgyowan/go-email-grpc/domain/email"

type Service interface {
	email.EmailService
}
