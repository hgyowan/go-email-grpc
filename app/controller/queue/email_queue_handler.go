package queue

import (
	"context"
	"encoding/json"

	"github.com/hgyowan/go-email-grpc/domain/email"
	"github.com/hgyowan/go-email-grpc/pkg/constant"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
)

func listenEmailQueueHandler(ctx context.Context, h *queueHandler) {
	eventCh, errorCh, _ := h.externalQueueListener.EmailQueueListener().Listen(ctx, constant.EmailSendQueueEventName)
	for {
		select {
		case events := <-eventCh:
			for _, event := range events {
				switch event.EventName {
				case constant.EmailSendQueueEventName:
					var param email.EmailServiceParam
					err := json.Unmarshal(event.Data, &param)
					if err != nil {
						pkgLogger.ZapLogger.Logger.Sugar().Error(pkgError.Wrap(err))
						continue
					}

					if err = h.service.SendTemplateEmail(ctx, param); err != nil {
						pkgLogger.ZapLogger.Logger.Sugar().Error(pkgError.Wrap(err))
					}

					if err = h.externalQueueListener.EmailQueueListener().DeleteMessage(ctx, event.ReceiptHandle); err != nil {
						pkgLogger.ZapLogger.Logger.Sugar().Error(pkgError.Wrap(err))
					}
				}
			}
		case err := <-errorCh:
			pkgLogger.ZapLogger.Logger.Sugar().Error(pkgError.Wrap(err))
		case <-ctx.Done():
			pkgLogger.ZapLogger.Logger.Info("Listener shutting down...")
			return
		}
	}
}
