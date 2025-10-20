package external

import (
	"github.com/hgyowan/go-email-grpc/domain"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgQueue "github.com/hgyowan/go-pkg-library/queue"
	"github.com/segmentio/kafka-go"
)

type externalQueueListener struct {
	emailQueueListener pkgQueue.EventListener
}

func (e *externalQueueListener) EmailQueueListener() pkgQueue.EventListener {
	return e.emailQueueListener
}

func MustNewExternalQueueListener() domain.ExternalQueueListener {
	return &externalQueueListener{emailQueueListener: pkgQueue.MustNewKafkaEventListener(kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{envs.EmailQueueBroker},
		GroupID:  envs.EmailQueueGroup,
		Topic:    envs.EmailQueueTopic,
		MinBytes: 1,
		MaxBytes: 1024 * 1024 * 1,
		MaxWait:  0,
	}), &pkgQueue.KafkaListenerConfig{
		ConsumerCount: 1,
	})}
}
