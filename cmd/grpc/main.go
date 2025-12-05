package main

import (
	"context"
	"github.com/hgyowan/go-email-grpc/app/controller/grpc"
	"github.com/hgyowan/go-email-grpc/app/controller/queue"
	"github.com/hgyowan/go-email-grpc/app/external"
	"github.com/hgyowan/go-email-grpc/app/repository"
	"github.com/hgyowan/go-email-grpc/app/service"
	pkgCrypto "github.com/hgyowan/go-pkg-library/crypto"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	pkgTrace "github.com/hgyowan/go-pkg-library/trace"
	"log"
)

func main() {
	pkgLogger.MustInitZapLogger()
	pkgCrypto.MustNewCryptoHelper([]byte(envs.MasterKey))

	if pkgLogger.ZapLogger == nil {
		log.Fatal("logger is nil")
	}

	gCtx, cancelFunc := context.WithCancel(context.Background())

	defer cancelFunc()

	if envs.ServiceType == envs.PrdType {
		shutdown := pkgTrace.InitTracer(gCtx, &pkgTrace.OpenTelemetryConfig{
			ServiceName: envs.ServerName,
			Endpoint:    envs.OpenTelemetryEndpoint,
		})
		defer shutdown()
	}

	grpcServer := external.MustNewGRPCServer()
	dbClient := external.MustNewExternalDB()
	repo := repository.NewRepository(dbClient)
	redisCli := external.MustNewExternalRedis()
	v := external.MustNewValidator()
	mailSender := external.MustNewEmailSenderV2("./internal/format/")
	queueListener := external.MustNewExternalQueueListener()
	queueEmitter := external.MustNewExternalQueueEmitter()
	svc := service.NewService(gCtx, repo, redisCli, queueListener, queueEmitter, mailSender, v)
	pkgLogger.ZapLogger.Logger.Info("Starting gRPC server on")

	q := queue.NewQueueHandler(svc, queueListener)
	go q.Listen(gCtx)

	grpcHandler := grpc.NewGRPCHandler(svc, grpcServer)
	grpcHandler.Listen(gCtx)

	pkgLogger.ZapLogger.Logger.Info("GRPC Server End")
}
