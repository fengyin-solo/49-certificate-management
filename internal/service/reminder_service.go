package service

import (
	"sort"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateReminder(input model.Reminder) (*model.Reminder, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCertificate(input.CertificateID); err != nil {
		return nil, model.NewValidationError("certificate_id", "证书不存在")
	}
	r := &model.Reminder{
		ID:            idgen.Hex(),
		CertificateID: input.CertificateID,
		Type:          input.Type,
		AdvanceDays:   input.AdvanceDays,
		TriggerAt:     input.TriggerAt,
		Status:        input.Status,
		Message:       input.Message,
	}
	if err := s.store.CreateReminder(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) ListReminders(filter model.ReminderFilter, page, size int) ([]*model.Reminder, int, error) {
	all := s.store.ListReminders()
	matched := make([]*model.Reminder, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].TriggerAt.Before(matched[j].TriggerAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Reminder{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetReminder(id string) (*model.Reminder, error) {
	return s.store.GetReminder(id)
}

func (s *Service) UpdateReminder(id string, input model.Reminder) (*model.Reminder, error) {
	r, err := s.store.GetReminder(id)
	if err != nil {
		return nil, err
	}
	if input.Type != "" {
		r.Type = input.Type
	}
	if input.AdvanceDays >= 0 {
		r.AdvanceDays = input.AdvanceDays
	}
	if !input.TriggerAt.IsZero() {
		r.TriggerAt = input.TriggerAt
	}
	if input.Status != "" {
		r.Status = input.Status
	}
	r.Message = input.Message
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateReminder(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteReminder(id string) error {
	return s.store.DeleteReminder(id)
}
