package email

import (
	"encoding/json"

	"github.com/hgyowan/go-email-grpc/pkg/constant"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgEmailV2 "github.com/hgyowan/go-pkg-library/mail/v2"
)

type TemplateMetaData interface {
	Unmarshal(b []byte) error
	Type() pkgEmailV2.EmailTemplateType
	GetSubject(langCode constant.LangCode, subject string) string
}

func NewEmailTemplateMetaData(templateType pkgEmailV2.EmailTemplateType) (TemplateMetaData, error) {
	switch templateType {
	case pkgEmailV2.EmailTemplateTypeVerifyEmail:
		return &VerifyEmail{}, nil
	case pkgEmailV2.EmailTemplateTypeJoinMessage:
		return &JoinMessage{}, nil
	case pkgEmailV2.EmailTemplateTypeJoinConfirm:
		return &JoinConfirm{}, nil
	case pkgEmailV2.EmailTemplateTypeInviteSend:
		return &InviteSend{}, nil
	}

	return nil, pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.WrongParam)
}

type VerifyEmail struct {
	VerifyCode string `json:"verifyCode"`
}

func (v *VerifyEmail) Unmarshal(b []byte) error {
	if err := json.Unmarshal(b, &v); err != nil {
		return pkgError.Wrap(err)
	}

	return nil
}

func (v *VerifyEmail) Type() pkgEmailV2.EmailTemplateType {
	return pkgEmailV2.EmailTemplateTypeVerifyEmail
}

func (v *VerifyEmail) GetSubject(langCode constant.LangCode, subject string) string {
	if subject != "" {
		return subject
	}

	return ""
}

type JoinMessage struct {
	WorkspaceName string `json:"workspaceName"`
	UserName      string `json:"userName"`
	UserEmail     string `json:"userEmail"`
	WorkspaceLink string `json:"workspaceLink"`
}

func (j *JoinMessage) Unmarshal(b []byte) error {
	if err := json.Unmarshal(b, &j); err != nil {
		return pkgError.Wrap(err)
	}

	return nil
}

func (j *JoinMessage) Type() pkgEmailV2.EmailTemplateType {
	return pkgEmailV2.EmailTemplateTypeJoinMessage
}

func (j *JoinMessage) GetSubject(langCode constant.LangCode, subject string) string {
	if subject != "" {
		return subject
	}

	return ""
}

type InviteSend struct {
	WorkspaceName string `json:"workspaceName"`
	JoinLink      string `json:"joinLink"`
}

func (i *InviteSend) Unmarshal(b []byte) error {
	if err := json.Unmarshal(b, &i); err != nil {
		return pkgError.Wrap(err)
	}

	return nil
}

func (i *InviteSend) Type() pkgEmailV2.EmailTemplateType {
	return pkgEmailV2.EmailTemplateTypeInviteSend
}

func (i *InviteSend) GetSubject(langCode constant.LangCode, subject string) string {
	if subject != "" {
		return subject
	}

	return ""
}

type JoinConfirm struct {
	WorkspaceName string `json:"workspaceName"`
	WorkspaceLink string `json:"workspaceLink"`
}

func (j *JoinConfirm) Unmarshal(b []byte) error {
	if err := json.Unmarshal(b, &j); err != nil {
		return pkgError.Wrap(err)
	}

	return nil
}

func (j *JoinConfirm) Type() pkgEmailV2.EmailTemplateType {
	return pkgEmailV2.EmailTemplateTypeJoinConfirm
}

func (j *JoinConfirm) GetSubject(langCode constant.LangCode, subject string) string {
	if subject != "" {
		return subject
	}

	return ""
}
