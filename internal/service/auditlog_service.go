package service

import (
	"sort"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateAuditLog(input model.AuditLog) (*model.AuditLog, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	a := &model.AuditLog{
		ID:         idgen.Hex(),
		Operator:   input.Operator,
		Action:     input.Action,
		TargetType: input.TargetType,
		TargetID:   input.TargetID,
		Detail:     input.Detail,
		CreatedAt:  input.CreatedAt,
	}
	if err := s.store.CreateAuditLog(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) ListAuditLogs(filter model.AuditLogFilter, page, size int) ([]*model.AuditLog, int, error) {
	all := s.store.ListAuditLogs()
	matched := make([]*model.AuditLog, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AuditLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetAuditLog(id string) (*model.AuditLog, error) {
	return s.store.GetAuditLog(id)
}

func (s *Service) DeleteAuditLog(id string) error {
	return s.store.DeleteAuditLog(id)
}
