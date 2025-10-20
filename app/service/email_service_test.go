package service

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/hgyowan/go-email-grpc/domain/email"
	"github.com/hgyowan/go-email-grpc/pkg/constant"
	pkgEmailV2 "github.com/hgyowan/go-pkg-library/mail/v2"
	"github.com/stretchr/testify/require"
)

func TestService_SendTemplateEmail(t *testing.T) {
	beforeEach()
	verifyEmail := email.VerifyEmail{VerifyCode: "123456"}
	b, _ := json.Marshal(verifyEmail)
	err := svc.TemplateEmailEmit(ctx, email.EmailServiceParam{
		List: []*email.RecipientRequest{
			{
				ID:               uuid.NewString(),
				LangCode:         constant.KO,
				TemplateType:     pkgEmailV2.EmailTemplateTypeVerifyEmail,
				ToEmails:         []string{"rydhkstptkd@naver.com"},
				Subject:          "test",
				TemplateMetaData: string(b),
			},
		},
	})
	require.NoError(t, err)
}
