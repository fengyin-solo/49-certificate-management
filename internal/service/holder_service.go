package service

import (
	"sort"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateHolder(input model.Holder) (*model.Holder, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	h := &model.Holder{
		ID:       idgen.Hex(),
		Name:     input.Name,
		IDNumber: input.IDNumber,
		Phone:    input.Phone,
		Email:    input.Email,
		Company:  input.Company,
		Status:   input.Status,
	}
	if err := s.store.CreateHolder(h); err != nil {
		return nil, err
	}
	return h, nil
}

func (s *Service) ListHolders(filter model.HolderFilter, page, size int) ([]*model.Holder, int, error) {
	all := s.store.ListHolders()
	matched := make([]*model.Holder, 0, len(all))
	for _, h := range all {
		if filter.Match(h) {
			matched = append(matched, h)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Name < matched[j].Name
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Holder{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetHolder(id string) (*model.Holder, error) {
	return s.store.GetHolder(id)
}

func (s *Service) UpdateHolder(id string, input model.Holder) (*model.Holder, error) {
	h, err := s.store.GetHolder(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		h.Name = input.Name
	}
	if input.IDNumber != "" {
		h.IDNumber = input.IDNumber
	}
	h.Phone = input.Phone
	h.Email = input.Email
	h.Company = input.Company
	if input.Status != "" {
		h.Status = input.Status
	}
	if err := h.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateHolder(h); err != nil {
		return nil, err
	}
	return h, nil
}

func (s *Service) DeleteHolder(id string) error {
	return s.store.DeleteHolder(id)
}
