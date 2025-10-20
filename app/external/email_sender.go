package external

import (
	"github.com/hgyowan/go-email-grpc/domain"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgEmailV2 "github.com/hgyowan/go-pkg-library/mail/v2"
	"time"
)

type externalMailSender struct {
	pkgEmailV2.EmailSenderV2
}

func (e *externalMailSender) Sender() pkgEmailV2.EmailSenderV2 {
	return e.EmailSenderV2
}

func MustNewEmailSenderV2(formatDirectory string) domain.ExternalEmailSender {
	emailSenderV2 := pkgEmailV2.MustNewEmailSender(&pkgEmailV2.EmailConfig{
		SMTPHost:         envs.SMTPServer,
		SMTPPort:         envs.SMTPPort,
		SMTPSender:       envs.SMTPSender,
		Username:         envs.SMTPAccount,
		Password:         envs.SMTPPassword,
		DelayBetweenMsg:  500 * time.Millisecond,
		MaxRetries:       3,
		EmailWorkerCount: 2,
	}, formatDirectory)

	return &externalMailSender{emailSenderV2}
}
