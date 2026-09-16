package model

import (
	"strings"
	"time"
)

const (
	ComplaintPending       = "pending"
	ComplaintInvestigating = "investigating"
	ComplaintResolved      = "resolved"
	ComplaintDismissed     = "dismissed"
)

var complaintTransitions = map[string]map[string]bool{
	ComplaintPending:       {ComplaintInvestigating: true},
	ComplaintInvestigating: {ComplaintResolved: true, ComplaintDismissed: true},
	ComplaintResolved:      {},
	ComplaintDismissed:     {},
}

func ComplaintCanTransition(from, to string) bool {
	if m, ok := complaintTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Complaint struct {
	ID           string    `json:"id"`
	CertificateID string    `json:"certificate_id"`
	Complainant  string    `json:"complainant"`
	Type         string    `json:"type"`
	Content      string    `json:"content"`
	Status       string    `json:"status"`
	Handler      string    `json:"handler"`
	HandleResult string    `json:"handle_result"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (c *Complaint) Validate() error {
	c.CertificateID = strings.TrimSpace(c.CertificateID)
	c.Complainant = strings.TrimSpace(c.Complainant)
	c.Type = strings.TrimSpace(c.Type)
	c.Content = strings.TrimSpace(c.Content)
	c.Handler = strings.TrimSpace(c.Handler)
	c.HandleResult = strings.TrimSpace(c.HandleResult)
	if c.CertificateID == "" {
		return NewValidationError("certificate_id", "证书 ID 不能为空")
	}
	if c.Complainant == "" {
		return NewValidationError("complainant", "投诉人不能为空")
	}
	if c.Type == "" {
		return NewValidationError("type", "投诉类型不能为空")
	}
	if c.Content == "" {
		return NewValidationError("content", "投诉内容不能为空")
	}
	if c.Status == "" {
		c.Status = ComplaintPending
	}
	if c.Status != ComplaintPending && c.Status != ComplaintInvestigating && c.Status != ComplaintResolved && c.Status != ComplaintDismissed {
		return NewValidationError("status", "投诉状态不合法")
	}
	return nil
}

type ComplaintFilter struct {
	CertificateID string
	Status        string
	Type          string
	Complainant   string
}

func (f ComplaintFilter) Match(c *Complaint) bool {
	if f.CertificateID != "" && c.CertificateID != f.CertificateID {
		return false
	}
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Type != "" && c.Type != f.Type {
		return false
	}
	if f.Complainant != "" && c.Complainant != f.Complainant {
		return false
	}
	return true
}
