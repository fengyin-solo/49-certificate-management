package model

import (
	"strings"
)

const (
	CategoryActive   = "active"
	CategoryDisabled = "disabled"
)

type Category struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ValidYears   int    `json:"valid_years"`
	RequireReview bool  `json:"require_review"`
	Status       string `json:"status"`
}

func (c *Category) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return NewValidationError("name", "分类名称不能为空")
	}
	if c.ValidYears < 0 {
		return NewValidationError("valid_years", "有效年限不能为负数")
	}
	if c.Status == "" {
		c.Status = CategoryActive
	}
	if c.Status != CategoryActive && c.Status != CategoryDisabled {
		return NewValidationError("status", "分类状态不合法")
	}
	return nil
}

type CategoryFilter struct {
	Status  string
	Keyword string
}

func (f CategoryFilter) Match(c *Category) bool {
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) {
			return false
		}
	}
	return true
}
