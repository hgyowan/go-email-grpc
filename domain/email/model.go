package email

import (
	"github.com/hgyowan/go-email-grpc/pkg/constant"
	pkgCrypto "github.com/hgyowan/go-pkg-library/crypto"
	pkgEmailV2 "github.com/hgyowan/go-pkg-library/mail/v2"
	"gorm.io/gorm"
	"time"
)

type EmailSendLog struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	EmailID      string    `gorm:"column:email_id"`
	Email        string    `gorm:"column:email" crypto:"type:fixed_cbc;context:email"`
	LangCode     string    `gorm:"column:lang_code"`
	TemplateType string    `gorm:"column:template_type"`
	MetaData     string    `gorm:"column:meta_data"`
	Status       string    `gorm:"column:status"`
	FailReason   string    `gorm:"column:fail_reason"`
	SendingAt    time.Time `gorm:"column:sending_at"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (esli *EmailSendLog) TableName() string {
	return "email_send_logs"
}

func (esli *EmailSendLog) BeforeCreate(*gorm.DB) error {
	return pkgCrypto.EncryptScheme(esli)
}

func (esli *EmailSendLog) AfterCreate(*gorm.DB) error {
	return pkgCrypto.DecryptScheme(esli)
}

func (esli *EmailSendLog) BeforeUpdate(*gorm.DB) error {
	return pkgCrypto.EncryptScheme(esli)
}

func (esli *EmailSendLog) AfterUpdate(*gorm.DB) error {
	return pkgCrypto.DecryptScheme(esli)
}

func (esli *EmailSendLog) AfterFind(*gorm.DB) error {
	return pkgCrypto.DecryptScheme(esli)
}

type Recipient struct {
	ToEmails         []string                     `json:"toEmails"`
	LangCode         constant.LangCode            `json:"langCode"`
	Subject          string                       `json:"subject"`
	TemplateType     pkgEmailV2.EmailTemplateType `json:"templateType"`
	TemplateMetaData interface{}                  `json:"templateMetaData"`
}
