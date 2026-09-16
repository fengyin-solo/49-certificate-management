package service

import (
	"sort"
	"time"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateNotification(input model.Notification) (*model.Notification, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetReminder(input.ReminderID); err != nil {
		return nil, model.NewValidationError("reminder_id", "提醒不存在")
	}
	n := &model.Notification{
		ID:         idgen.Hex(),
		ReminderID: input.ReminderID,
		Recipient:  input.Recipient,
		Channel:    input.Channel,
		Status:     input.Status,
		SentAt:     input.SentAt,
		Content:    input.Content,
	}
	if err := s.store.CreateNotification(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) ListNotifications(filter model.NotificationFilter, page, size int) ([]*model.Notification, int, error) {
	all := s.store.ListNotifications()
	matched := make([]*model.Notification, 0, len(all))
	for _, n := range all {
		if filter.Match(n) {
			matched = append(matched, n)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].SentAt.After(matched[j].SentAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Notification{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetNotification(id string) (*model.Notification, error) {
	return s.store.GetNotification(id)
}

func (s *Service) SendNotification(id string) (*model.Notification, error) {
	n, err := s.store.GetNotification(id)
	if err != nil {
		return nil, err
	}
	n.Status = model.NotificationSent
	n.SentAt = time.Now()
	if err := s.store.UpdateNotification(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) DeleteNotification(id string) error {
	return s.store.DeleteNotification(id)
}
