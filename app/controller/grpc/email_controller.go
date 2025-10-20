package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/hgyowan/go-email-grpc/domain/email"
	emailV1 "github.com/hgyowan/go-email-grpc/gen/email/v1"
	"github.com/hgyowan/go-email-grpc/internal"
	"github.com/hgyowan/go-email-grpc/pkg/constant"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgEmailV2 "github.com/hgyowan/go-pkg-library/mail/v2"
)

func registerEmailGRPCHandler(h *grpcHandler) {
	h.EmailServiceServer = &emailGRPCHandler{h: h}
	emailV1.RegisterEmailServiceServer(h.externalGRPCServer.Server(), h)
}

type emailGRPCHandler struct {
	h *grpcHandler
}

func (e *emailGRPCHandler) SendTemplateEmail(ctx context.Context, request *emailV1.SendTemplateEmailRequest) (*emailV1.SendTemplateEmailResponse, error) {
	list := make([]*email.RecipientRequest, 0, len(request.GetList()))
	for _, item := range request.GetList() {
		if item.GetTemplateMetadata() == "" {
			return nil, pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.WrongParam)
		}

		metaData, err := internal.ParseMetadata([]byte(item.GetTemplateMetadata()))
		if err != nil {
			return nil, pkgError.WrapWithCode(err, pkgError.WrongParam)
		}

		list = append(list, &email.RecipientRequest{
			LangCode:         constant.LangCode(item.GetLangCode()),
			TemplateType:     pkgEmailV2.EmailTemplateType(item.GetTemplateType()),
			ToEmails:         item.GetToEmails(),
			Subject:          item.GetSubject(),
			TemplateMetaData: metaData,
			ID:               uuid.NewString(),
		})
	}

	if err := e.h.service.TemplateEmailEmit(ctx, email.EmailServiceParam{List: list}); err != nil {
		return nil, pkgError.Wrap(err)
	}

	return &emailV1.SendTemplateEmailResponse{}, nil
}
