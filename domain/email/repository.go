package email

import "context"

type EmailRepository interface {
	CreateEmailSendLogBatch(ctx context.Context, param []*EmailSendLog) error
}
