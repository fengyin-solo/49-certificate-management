package model

import (
	"strings"
)

const (
	HolderActive   = "active"
	HolderDisabled = "disabled"
)

type Holder struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IDNumber string `json:"id_number"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Company string `json:"company"`
	Status  string `json:"status"`
}

func (h *Holder) Validate() error {
	h.Name = strings.TrimSpace(h.Name)
	h.IDNumber = strings.TrimSpace(h.IDNumber)
	h.Phone = strings.TrimSpace(h.Phone)
	h.Email = strings.TrimSpace(h.Email)
	h.Company = strings.TrimSpace(h.Company)
	if h.Name == "" {
		return NewValidationError("name", "持有人姓名不能为空")
	}
	if h.IDNumber == "" {
		return NewValidationError("id_number", "证件号不能为空")
	}
	if h.Status == "" {
		h.Status = HolderActive
	}
	if h.Status != HolderActive && h.Status != HolderDisabled {
		return NewValidationError("status", "持有人状态不合法")
	}
	return nil
}

type HolderFilter struct {
	Status  string
	Keyword string
}

func (f HolderFilter) Match(h *Holder) bool {
	if f.Status != "" && h.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(h.Name), k) && !strings.Contains(strings.ToLower(h.Company), k) {
			return false
		}
	}
	return true
}
