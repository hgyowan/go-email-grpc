package email

import "context"

type EmailService interface {
	SendTemplateEmail(ctx context.Context, param EmailServiceParam) error
	TemplateEmailEmit(ctx context.Context, param EmailServiceParam) error
}
