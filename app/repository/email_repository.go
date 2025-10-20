package repository

import (
	"github.com/hgyowan/go-email-grpc/domain/email"
	pkgError "github.com/hgyowan/go-pkg-library/error"
)

func registerEmailRepository(r *repository) {
	r.EmailRepository = &emailRepository{repository: r}
}

type emailRepository struct {
	repository *repository
}

func (e *emailRepository) CreateEmailSendLogBatch(param []*email.EmailSendLog) error {
	return pkgError.Wrap(e.repository.externalGormClient.DB().CreateInBatches(&param, 500).Error)
}
