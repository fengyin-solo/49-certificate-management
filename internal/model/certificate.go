package model

import (
	"strings"
	"time"
)

const (
	CertificateActive   = "active"
	CertificateExpired  = "expired"
	CertificateRevoked  = "revoked"
	CertificateSuspended = "suspended"
)

var certificateTransitions = map[string]map[string]bool{
	CertificateActive:   {CertificateExpired: true, CertificateRevoked: true, CertificateSuspended: true},
	CertificateSuspended: {CertificateActive: true, CertificateRevoked: true},
	CertificateExpired:  {},
	CertificateRevoked:  {},
}

func CertificateCanTransition(from, to string) bool {
	if m, ok := certificateTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Certificate struct {
	ID         string    `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	CategoryID string    `json:"category_id"`
	HolderID   string    `json:"holder_id"`
	IssuerID   string    `json:"issuer_id"`
	IssueDate  string    `json:"issue_date"`
	ExpireDate string    `json:"expire_date"`
	Status     string    `json:"status"`
	Remarks    string    `json:"remarks"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *Certificate) Validate() error {
	c.Code = strings.TrimSpace(c.Code)
	c.Name = strings.TrimSpace(c.Name)
	c.CategoryID = strings.TrimSpace(c.CategoryID)
	c.HolderID = strings.TrimSpace(c.HolderID)
	c.IssuerID = strings.TrimSpace(c.IssuerID)
	c.IssueDate = strings.TrimSpace(c.IssueDate)
	c.ExpireDate = strings.TrimSpace(c.ExpireDate)
	if c.Code == "" {
		return NewValidationError("code", "证书编号不能为空")
	}
	if c.Name == "" {
		return NewValidationError("name", "证书名称不能为空")
	}
	if c.CategoryID == "" {
		return NewValidationError("category_id", "分类不能为空")
	}
	if c.HolderID == "" {
		return NewValidationError("holder_id", "持有人不能为空")
	}
	if c.IssuerID == "" {
		return NewValidationError("issuer_id", "发证机构不能为空")
	}
	if c.IssueDate == "" {
		return NewValidationError("issue_date", "签发日期不能为空")
	}
	if c.ExpireDate == "" {
		return NewValidationError("expire_date", "到期日期不能为空")
	}
	if c.Status == "" {
		c.Status = CertificateActive
	}
	if c.Status != CertificateActive && c.Status != CertificateExpired && c.Status != CertificateRevoked && c.Status != CertificateSuspended {
		return NewValidationError("status", "证书状态不合法")
	}
	return nil
}

type CertificateFilter struct {
	CategoryID string
	HolderID   string
	IssuerID   string
	Status     string
	Keyword    string
}

func (f CertificateFilter) Match(c *Certificate) bool {
	if f.CategoryID != "" && c.CategoryID != f.CategoryID {
		return false
	}
	if f.HolderID != "" && c.HolderID != f.HolderID {
		return false
	}
	if f.IssuerID != "" && c.IssuerID != f.IssuerID {
		return false
	}
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) && !strings.Contains(strings.ToLower(c.Code), k) {
			return false
		}
	}
	return true
}
