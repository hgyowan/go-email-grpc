package external

import (
	"github.com/hgyowan/go-email-grpc/domain"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgQueue "github.com/hgyowan/go-pkg-library/queue"
	"github.com/segmentio/kafka-go"
)

type externalQueueEmitter struct {
	emailEmitter pkgQueue.EventEmitter
}

func (e *externalQueueEmitter) EmailQueueEmitter() pkgQueue.EventEmitter {
	return e.emailEmitter
}

func MustNewExternalQueueEmitter() domain.ExternalQueueEmitter {
	return &externalQueueEmitter{emailEmitter: pkgQueue.MustNewKafkaEventEmitter(&kafka.Writer{
		Addr:     kafka.TCP(envs.EmailQueueBroker),
		Topic:    envs.EmailQueueTopic,
		Balancer: kafka.CRC32Balancer{},
	})}
}
