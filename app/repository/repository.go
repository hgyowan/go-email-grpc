package repository

import (
	"context"
	"github.com/hgyowan/go-email-grpc/domain"
	"github.com/hgyowan/go-email-grpc/domain/email"
)

type repository struct {
	email.EmailRepository
	externalGormClient domain.ExternalDBClient
}

func NewRepository(externalGormClient domain.ExternalDBClient) domain.Repository {
	r := &repository{
		externalGormClient: externalGormClient,
	}
	r.register()

	return r
}

func (r *repository) register() {
	registerEmailRepository(r)
}

func (r *repository) WithTransaction(ctx context.Context, fn func(txRepo domain.Repository) error) error {
	tx := r.externalGormClient.DB().WithContext(ctx).Begin() // 트랜잭션 시작
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rc := recover(); rc != nil {
			tx.Rollback()
			panic(rc)
		}
	}()

	txRepo := &repository{
		externalGormClient: r.externalGormClient.NewTxDB(tx),
	}

	txRepo.register()

	err := fn(txRepo)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
