package model

import (
	"strings"
	"time"
)

type AuditLog struct {
	ID         string    `json:"id"`
	Operator   string    `json:"operator"`
	Action     string    `json:"action"`
	TargetType string    `json:"target_type"`
	TargetID   string    `json:"target_id"`
	Detail     string    `json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}

func (a *AuditLog) Validate() error {
	a.Operator = strings.TrimSpace(a.Operator)
	a.Action = strings.TrimSpace(a.Action)
	a.TargetType = strings.TrimSpace(a.TargetType)
	a.TargetID = strings.TrimSpace(a.TargetID)
	if a.Operator == "" {
		return NewValidationError("operator", "操作人不能为空")
	}
	if a.Action == "" {
		return NewValidationError("action", "操作行为不能为空")
	}
	if a.TargetType == "" {
		return NewValidationError("target_type", "目标类型不能为空")
	}
	if a.TargetID == "" {
		return NewValidationError("target_id", "目标 ID 不能为空")
	}
	return nil
}

type AuditLogFilter struct {
	Operator   string
	Action     string
	TargetType string
	TargetID   string
}

func (f AuditLogFilter) Match(a *AuditLog) bool {
	if f.Operator != "" && a.Operator != f.Operator {
		return false
	}
	if f.Action != "" && a.Action != f.Action {
		return false
	}
	if f.TargetType != "" && a.TargetType != f.TargetType {
		return false
	}
	if f.TargetID != "" && a.TargetID != f.TargetID {
		return false
	}
	return true
}
