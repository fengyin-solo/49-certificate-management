package model

import (
	"strings"
	"time"
)

const (
	RenewalPending   = "pending"
	RenewalApproved  = "approved"
	RenewalRejected  = "rejected"
	RenewalCompleted = "completed"
)

var renewalTransitions = map[string]map[string]bool{
	RenewalPending:  {RenewalApproved: true, RenewalRejected: true},
	RenewalApproved: {RenewalCompleted: true},
	RenewalRejected: {},
	RenewalCompleted: {},
}

func RenewalCanTransition(from, to string) bool {
	if m, ok := renewalTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Renewal struct {
	ID             string    `json:"id"`
	CertificateID  string    `json:"certificate_id"`
	Applicant      string    `json:"applicant"`
	Reason         string    `json:"reason"`
	Status         string    `json:"status"`
	ReviewComment  string    `json:"review_comment"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (r *Renewal) Validate() error {
	r.CertificateID = strings.TrimSpace(r.CertificateID)
	r.Applicant = strings.TrimSpace(r.Applicant)
	r.Reason = strings.TrimSpace(r.Reason)
	if r.CertificateID == "" {
		return NewValidationError("certificate_id", "证书 ID 不能为空")
	}
	if r.Applicant == "" {
		return NewValidationError("applicant", "申请人不能为空")
	}
	if r.Status == "" {
		r.Status = RenewalPending
	}
	if r.Status != RenewalPending && r.Status != RenewalApproved && r.Status != RenewalRejected && r.Status != RenewalCompleted {
		return NewValidationError("status", "续期状态不合法")
	}
	return nil
}

type RenewalFilter struct {
	CertificateID string
	Status        string
	Applicant     string
}

func (f RenewalFilter) Match(r *Renewal) bool {
	if f.CertificateID != "" && r.CertificateID != f.CertificateID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Applicant != "" && r.Applicant != f.Applicant {
		return false
	}
	return true
}
