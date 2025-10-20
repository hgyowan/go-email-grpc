package service

import (
	"context"

	"github.com/hgyowan/go-email-grpc/app/external"
	"github.com/hgyowan/go-email-grpc/app/repository"
	"github.com/hgyowan/go-email-grpc/domain"
	pkgCrypto "github.com/hgyowan/go-pkg-library/crypto"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
)

var ctx context.Context
var svc domain.Service

func beforeEach() {
	pkgLogger.MustInitZapLogger()
	pkgCrypto.MustNewCryptoHelper([]byte(envs.MasterKey))
	ctx = context.Background()
	db := external.MustNewExternalDB()
	redis := external.MustNewExternalRedis()
	v := external.MustNewValidator()
	mailSender := external.MustNewEmailSenderV2("/Users/hwang-gyowan/go/src/church-financial-core-grpc/internal/format/")
	repo := repository.NewRepository(db)
	queueListener := external.MustNewExternalQueueListener()
	queueEmitter := external.MustNewExternalQueueEmitter()
	svc = NewService(ctx, repo, redis, queueListener, queueEmitter, mailSender, v)
}
