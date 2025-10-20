package domain

import (
	"github.com/go-playground/validator/v10"
	pkgGrpc "github.com/hgyowan/go-pkg-library/grpc-library/grpc"
	pkgEmailV2 "github.com/hgyowan/go-pkg-library/mail/v2"
	pkgQueue "github.com/hgyowan/go-pkg-library/queue"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ExternalGRPCServer interface {
	Server() pkgGrpc.GrpcServer
	Port() string
}

type ExternalDBClient interface {
	DB() *gorm.DB
	NewTxDB(tx *gorm.DB) ExternalDBClient
}

type ExternalRedisClient interface {
	Redis() redis.Cmdable
}

type ExternalValidator interface {
	Validator() *validator.Validate
}

type ExternalQueueListener interface {
	EmailQueueListener() pkgQueue.EventListener
}

type ExternalQueueEmitter interface {
	EmailQueueEmitter() pkgQueue.EventEmitter
}

type ExternalEmailSender interface {
	Sender() pkgEmailV2.EmailSenderV2
}
