package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/hgyowan/go-email-grpc/domain/email"
	"github.com/hgyowan/go-email-grpc/pkg/constant"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	pkgEmailV2 "github.com/hgyowan/go-pkg-library/mail/v2"
	pkgQueue "github.com/hgyowan/go-pkg-library/queue"
	"github.com/samber/lo"
)

func registerEmailService(ctx context.Context, s *service) {
	duplicateChecker := pkgEmailV2.WithDuplicateChecker(func(recipient *pkgEmailV2.Recipient) bool {
		ok, _ := s.externalRedisClient.Redis().SetNX(ctx, fmt.Sprintf("emailSend:%s", recipient.ID), "sent", 24*time.Hour*3).Result()
		if !ok {
			// true: 중복 메일
			return true
		}

		return false
	})

	templateFunc := pkgEmailV2.WithTemplateFuncMap(template.FuncMap{
		"minusOne": func(x int) int { return x - 1 },
	})

	s.externalEmailSender.Sender().SetOptions(duplicateChecker, templateFunc)

	s.EmailService = &emailService{service: s}
}

type emailService struct {
	service *service
}

func (e *emailService) SendTemplateEmail(ctx context.Context, param email.EmailServiceParam) error {
	responseCh := e.service.externalEmailSender.Sender().SendMailWithTemplateV2Parallel(lo.FilterMap(param.List, func(item *email.RecipientRequest, index int) (*pkgEmailV2.Recipient, bool) {
		tmpl, err := email.NewEmailTemplateMetaData(item.TemplateType)
		if err != nil {
			pkgLogger.ZapLogger.Logger.Sugar().Error(err)
			return nil, false
		}

		if err = tmpl.Unmarshal([]byte(item.TemplateMetaData)); err != nil {
			pkgLogger.ZapLogger.Logger.Sugar().Error(err)
			return nil, false
		}

		if strings.ToUpper(string(item.LangCode)) != "KO" {
			item.LangCode = "EN"
		}

		return &pkgEmailV2.Recipient{
			ID:               item.ID,
			LangCode:         string(item.LangCode),
			TemplateType:     item.TemplateType,
			ToEmails:         item.ToEmails,
			Subject:          tmpl.GetSubject(item.LangCode, item.Subject),
			TemplateMetaData: tmpl,
		}, true
	}))

	for response := range responseCh {
		if response.Status == "DUPLICATE" {
			continue
		}

		mails := make([]*email.EmailSendLog, 0, len(response.Emails))
		for _, mail := range response.Emails {
			mails = append(mails, &email.EmailSendLog{
				EmailID:      response.MailID,
				Email:        mail,
				LangCode:     response.LangCode,
				TemplateType: string(response.TemplateType),
				MetaData:     response.TemplateData,
				Status:       response.Status,
				FailReason:   string(response.FailReason),
				SendingAt:    response.SendingDate,
				CreatedAt:    time.Now().UTC(),
			})
		}

		if err := e.service.repo.CreateEmailSendLogBatch(mails); err != nil {
			pkgLogger.ZapLogger.Logger.Sugar().Error(err)
		}
	}

	return nil
}

func (e *emailService) TemplateEmailEmit(ctx context.Context, param email.EmailServiceParam) error {
	if len(param.List) == 0 {
		return nil
	}

	const maxMessageSize = 1024 * 1024 // 1024KB

	var message []*email.RecipientRequest
	batchSize := 0

	pushMessage := func(message []*email.RecipientRequest) error {
		if len(message) == 0 {
			return nil
		}

		data, err := json.Marshal(email.EmailServiceParam{List: message})
		if err != nil {
			return pkgError.WrapWithCode(err, pkgError.Get)
		}

		if err := e.service.externalQueueEmitter.EmailQueueEmitter().Emit(ctx, pkgQueue.Event{
			EventName:      constant.EmailSendQueueEventName,
			MessageGroupId: constant.EmailSendQueueEventName,
			Data:           data,
		}); err != nil {
			return pkgError.WrapWithCode(err, pkgError.Create)
		}

		return nil
	}

	for _, recipient := range param.List {
		recipientData, err := json.Marshal(recipient)
		if err != nil {
			return pkgError.WrapWithCode(err, pkgError.Get)
		}

		// 현재 배치에 추가했을 때 maxMessageSize 초과 시 flush
		if batchSize+len(recipientData) > maxMessageSize {
			if err := pushMessage(message); err != nil {
				return err
			}
			message = nil
			batchSize = 0
		}

		message = append(message, recipient)
		batchSize += len(recipientData)
	}

	// 마지막 남은 message push
	if err := pushMessage(message); err != nil {
		return err
	}

	return nil
}
