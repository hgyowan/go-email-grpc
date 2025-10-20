package email

import pkgEmailV2 "github.com/hgyowan/go-pkg-library/mail/v2"

type GetEmailTemplateSampleByTemplateTypeAndLangCodeDBParam struct {
	TemplateType pkgEmailV2.EmailTemplateType `json:"templateType"`
	LangCode     string                       `json:"langCode"`
}
