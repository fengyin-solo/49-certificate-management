package model

import (
	"strings"
	"time"
)

const (
	ReviewPending = "pending"
	ReviewPassed  = "passed"
	ReviewFailed  = "failed"
)

var reviewTransitions = map[string]map[string]bool{
	ReviewPending: {ReviewPassed: true, ReviewFailed: true},
	ReviewPassed:  {},
	ReviewFailed:  {},
}

func ReviewCanTransition(from, to string) bool {
	if m, ok := reviewTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Review struct {
	ID             string    `json:"id"`
	CertificateID  string    `json:"certificate_id"`
	Cycle          int       `json:"cycle"`
	ReviewDate     string    `json:"review_date"`
	Result         string    `json:"result"`
	NextReviewDate string    `json:"next_review_date"`
	Reviewer       string    `json:"reviewer"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (r *Review) Validate() error {
	r.CertificateID = strings.TrimSpace(r.CertificateID)
	r.ReviewDate = strings.TrimSpace(r.ReviewDate)
	r.Result = strings.TrimSpace(r.Result)
	r.NextReviewDate = strings.TrimSpace(r.NextReviewDate)
	r.Reviewer = strings.TrimSpace(r.Reviewer)
	if r.CertificateID == "" {
		return NewValidationError("certificate_id", "证书 ID 不能为空")
	}
	if r.Cycle <= 0 {
		return NewValidationError("cycle", "复审周期必须大于 0")
	}
	if r.ReviewDate == "" {
		return NewValidationError("review_date", "复审日期不能为空")
	}
	if r.Reviewer == "" {
		return NewValidationError("reviewer", "复审人不能为空")
	}
	if r.Result == "" {
		r.Result = ReviewFailed
	}
	if r.Result != ReviewPassed && r.Result != ReviewFailed {
		return NewValidationError("result", "复审结果不合法")
	}
	if r.Status == "" {
		r.Status = ReviewPending
	}
	if r.Status != ReviewPending && r.Status != ReviewPassed && r.Status != ReviewFailed {
		return NewValidationError("status", "复审状态不合法")
	}
	return nil
}

type ReviewFilter struct {
	CertificateID string
	Status        string
	Reviewer      string
}

func (f ReviewFilter) Match(r *Review) bool {
	if f.CertificateID != "" && r.CertificateID != f.CertificateID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Reviewer != "" && r.Reviewer != f.Reviewer {
		return false
	}
	return true
}
