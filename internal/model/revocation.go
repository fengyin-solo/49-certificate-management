package model

import (
	"strings"
	"time"
)

type Revocation struct {
	ID            string    `json:"id"`
	CertificateID string    `json:"certificate_id"`
	Reason        string    `json:"reason"`
	Operator      string    `json:"operator"`
	RevokedAt     time.Time `json:"revoked_at"`
	Note          string    `json:"note"`
}

func (r *Revocation) Validate() error {
	r.CertificateID = strings.TrimSpace(r.CertificateID)
	r.Reason = strings.TrimSpace(r.Reason)
	r.Operator = strings.TrimSpace(r.Operator)
	if r.CertificateID == "" {
		return NewValidationError("certificate_id", "证书 ID 不能为空")
	}
	if r.Reason == "" {
		return NewValidationError("reason", "吊销原因不能为空")
	}
	if r.Operator == "" {
		return NewValidationError("operator", "操作人不能为空")
	}
	return nil
}

type RevocationFilter struct {
	CertificateID string
	Operator      string
}

func (f RevocationFilter) Match(r *Revocation) bool {
	if f.CertificateID != "" && r.CertificateID != f.CertificateID {
		return false
	}
	if f.Operator != "" && r.Operator != f.Operator {
		return false
	}
	return true
}
