package model

import (
	"strings"
	"time"
)

const (
	VerificationOnline  = "online"
	VerificationOffline = "offline"
	VerificationManual  = "manual"
)

const (
	VerifyValid   = "valid"
	VerifyInvalid = "invalid"
	VerifyExpired = "expired"
)

type Verification struct {
	ID            string    `json:"id"`
	CertificateID string    `json:"certificate_id"`
	Method        string    `json:"method"`
	Result        string    `json:"result"`
	Verifier      string    `json:"verifier"`
	VerifiedAt    time.Time `json:"verified_at"`
	Note          string    `json:"note"`
}

func (v *Verification) Validate() error {
	v.CertificateID = strings.TrimSpace(v.CertificateID)
	v.Method = strings.TrimSpace(v.Method)
	v.Result = strings.TrimSpace(v.Result)
	v.Verifier = strings.TrimSpace(v.Verifier)
	if v.CertificateID == "" {
		return NewValidationError("certificate_id", "证书 ID 不能为空")
	}
	if v.Method == "" {
		return NewValidationError("method", "验真方式不能为空")
	}
	if v.Method != VerificationOnline && v.Method != VerificationOffline && v.Method != VerificationManual {
		return NewValidationError("method", "验真方式不合法")
	}
	if v.Result == "" {
		return NewValidationError("result", "验真结果不能为空")
	}
	if v.Result != VerifyValid && v.Result != VerifyInvalid && v.Result != VerifyExpired {
		return NewValidationError("result", "验真结果不合法")
	}
	if v.Verifier == "" {
		return NewValidationError("verifier", "验真人不能为空")
	}
	return nil
}

type VerificationFilter struct {
	CertificateID string
	Method        string
	Result        string
}

func (f VerificationFilter) Match(v *Verification) bool {
	if f.CertificateID != "" && v.CertificateID != f.CertificateID {
		return false
	}
	if f.Method != "" && v.Method != f.Method {
		return false
	}
	if f.Result != "" && v.Result != f.Result {
		return false
	}
	return true
}
