package model

import (
	"strings"
	"time"
)

const (
	ChannelEmail = "email"
	ChannelSMS   = "sms"
	ChannelInApp = "inapp"
)

const (
	NotificationPending = "pending"
	NotificationSent    = "sent"
	NotificationFailed  = "failed"
)

type Notification struct {
	ID         string    `json:"id"`
	ReminderID string    `json:"reminder_id"`
	Recipient  string    `json:"recipient"`
	Channel    string    `json:"channel"`
	Status     string    `json:"status"`
	SentAt     time.Time `json:"sent_at"`
	Content    string    `json:"content"`
}

func (n *Notification) Validate() error {
	n.ReminderID = strings.TrimSpace(n.ReminderID)
	n.Recipient = strings.TrimSpace(n.Recipient)
	n.Channel = strings.TrimSpace(n.Channel)
	n.Content = strings.TrimSpace(n.Content)
	if n.ReminderID == "" {
		return NewValidationError("reminder_id", "提醒 ID 不能为空")
	}
	if n.Recipient == "" {
		return NewValidationError("recipient", "接收人不能为空")
	}
	if n.Channel == "" {
		return NewValidationError("channel", "通知渠道不能为空")
	}
	if n.Channel != ChannelEmail && n.Channel != ChannelSMS && n.Channel != ChannelInApp {
		return NewValidationError("channel", "通知渠道不合法")
	}
	if n.Status == "" {
		n.Status = NotificationPending
	}
	if n.Status != NotificationPending && n.Status != NotificationSent && n.Status != NotificationFailed {
		return NewValidationError("status", "通知状态不合法")
	}
	return nil
}

type NotificationFilter struct {
	ReminderID string
	Channel    string
	Status     string
}

func (f NotificationFilter) Match(n *Notification) bool {
	if f.ReminderID != "" && n.ReminderID != f.ReminderID {
		return false
	}
	if f.Channel != "" && n.Channel != f.Channel {
		return false
	}
	if f.Status != "" && n.Status != f.Status {
		return false
	}
	return true
}
