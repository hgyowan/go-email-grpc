package email

import (
	"github.com/hgyowan/go-email-grpc/pkg/constant"
	pkgEmailV2 "github.com/hgyowan/go-pkg-library/mail/v2"
)

type RecipientRequest struct {
	ID               string                       `json:"id"`
	LangCode         constant.LangCode            `json:"langCode"`
	TemplateType     pkgEmailV2.EmailTemplateType `json:"templateType"`
	ToEmails         []string                     `json:"toEmails"` // 변하지 않는 템플릿 메일일경우 여러명에게 전송 가능
	Subject          string                       `json:"subject"`
	TemplateMetaData string                       `json:"templateMetaData"` // Meta Data JsonString
}

type EmailServiceParam struct {
	List []*RecipientRequest `json:"list"`
}

type GetTemplateSampleServiceParam struct {
	TemplateType pkgEmailV2.EmailTemplateType `json:"templateType"`
	LangCode     string                       `json:"langCode"`
}

type GetTemplateSampleResponse struct {
	SampleHTML string `json:"sampleHTML"`
}
