package model

import (
	"strings"
	"time"
)

const (
	ReminderBeforeExpire = "before_expire"
	ReminderReviewDue    = "review_due"
)

const (
	ReminderPending = "pending"
	ReminderSent    = "sent"
	ReminderClosed  = "closed"
)

type Reminder struct {
	ID            string    `json:"id"`
	CertificateID string    `json:"certificate_id"`
	Type          string    `json:"type"`
	AdvanceDays   int       `json:"advance_days"`
	TriggerAt     time.Time `json:"trigger_at"`
	Status        string    `json:"status"`
	Message       string    `json:"message"`
}

func (r *Reminder) Validate() error {
	r.CertificateID = strings.TrimSpace(r.CertificateID)
	r.Type = strings.TrimSpace(r.Type)
	r.Message = strings.TrimSpace(r.Message)
	if r.CertificateID == "" {
		return NewValidationError("certificate_id", "证书 ID 不能为空")
	}
	if r.Type == "" {
		return NewValidationError("type", "提醒类型不能为空")
	}
	if r.Type != ReminderBeforeExpire && r.Type != ReminderReviewDue {
		return NewValidationError("type", "提醒类型不合法")
	}
	if r.AdvanceDays < 0 {
		return NewValidationError("advance_days", "提前天数不能为负数")
	}
	if r.Status == "" {
		r.Status = ReminderPending
	}
	if r.Status != ReminderPending && r.Status != ReminderSent && r.Status != ReminderClosed {
		return NewValidationError("status", "提醒状态不合法")
	}
	return nil
}

type ReminderFilter struct {
	CertificateID string
	Type          string
	Status        string
}

func (f ReminderFilter) Match(r *Reminder) bool {
	if f.CertificateID != "" && r.CertificateID != f.CertificateID {
		return false
	}
	if f.Type != "" && r.Type != f.Type {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	return true
}
