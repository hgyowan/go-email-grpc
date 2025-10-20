package main

import (
	"context"
	"github.com/hgyowan/go-email-grpc/app/controller/grpc"
	"github.com/hgyowan/go-email-grpc/app/external"
	"github.com/hgyowan/go-email-grpc/app/repository"
	"github.com/hgyowan/go-email-grpc/app/service"
	"github.com/hgyowan/go-email-grpc/cmd/queue"
	pkgCrypto "github.com/hgyowan/go-pkg-library/crypto"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	pkgLogger.MustInitZapLogger()
	pkgCrypto.MustNewCryptoHelper([]byte(envs.MasterKey))

	if pkgLogger.ZapLogger == nil {
		log.Fatal("logger is nil")
	}

	bCtx, cancelFunc := context.WithCancel(context.Background())
	group, gCtx := errgroup.WithContext(bCtx)
	doneChan := make(chan struct{}, 1)
	grpcServer := external.MustNewGRPCServer()
	dbClient := external.MustNewExternalDB()
	repo := repository.NewRepository(dbClient)
	redisCli := external.MustNewExternalRedis()
	v := external.MustNewValidator()
	mailSender := external.MustNewEmailSenderV2("./internal/format/")
	queueListener := external.MustNewExternalQueueListener()
	queueEmitter := external.MustNewExternalQueueEmitter()
	svc := service.NewService(gCtx, repo, redisCli, queueListener, queueEmitter, mailSender, v)
	q := queue.NewQueueHandler(svc, queueListener)
	pkgLogger.ZapLogger.Logger.Info("Starting gRPC server on")

	group.Go(func() error {
		q.Listen(gCtx)
		pkgLogger.ZapLogger.Logger.Fatal("Queue Listen Handler End")
		doneChan <- struct{}{}
		return nil
	})

	group.Go(func() error {
		grpc.NewGRPCHandler(svc, grpcServer).Listen(gCtx)
		pkgLogger.ZapLogger.Logger.Fatal("GRPC Handler End")
		doneChan <- struct{}{}
		return nil
	})

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer close(interrupt)

	select {
	case <-doneChan:
		cancelFunc()
	case <-interrupt:
		cancelFunc()
	}

	if err := group.Wait(); err != nil {
		pkgLogger.ZapLogger.Logger.Fatal(err.Error())
	}

	pkgLogger.ZapLogger.Logger.Info("GRPC Server End")
}
