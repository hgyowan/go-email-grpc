package email

type EmailRepository interface {
	CreateEmailSendLogBatch(param []*EmailSendLog) error
}
