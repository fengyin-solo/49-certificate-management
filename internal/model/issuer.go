package model

import (
	"strings"
)

const (
	IssuerActive   = "active"
	IssuerDisabled = "disabled"
)

type Issuer struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Level        string `json:"level"`
	ContactPhone string `json:"contact_phone"`
	Address      string `json:"address"`
	Status       string `json:"status"`
}

func (i *Issuer) Validate() error {
	i.Name = strings.TrimSpace(i.Name)
	i.Level = strings.TrimSpace(i.Level)
	i.ContactPhone = strings.TrimSpace(i.ContactPhone)
	i.Address = strings.TrimSpace(i.Address)
	if i.Name == "" {
		return NewValidationError("name", "机构名称不能为空")
	}
	if i.Level == "" {
		return NewValidationError("level", "机构级别不能为空")
	}
	if i.Status == "" {
		i.Status = IssuerActive
	}
	if i.Status != IssuerActive && i.Status != IssuerDisabled {
		return NewValidationError("status", "机构状态不合法")
	}
	return nil
}

type IssuerFilter struct {
	Status  string
	Keyword string
}

func (f IssuerFilter) Match(i *Issuer) bool {
	if f.Status != "" && i.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(i.Name), k) {
			return false
		}
	}
	return true
}
