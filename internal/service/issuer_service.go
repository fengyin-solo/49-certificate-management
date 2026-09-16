package service

import (
	"sort"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateIssuer(input model.Issuer) (*model.Issuer, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	i := &model.Issuer{
		ID:           idgen.Hex(),
		Name:         input.Name,
		Level:        input.Level,
		ContactPhone: input.ContactPhone,
		Address:      input.Address,
		Status:       input.Status,
	}
	if err := s.store.CreateIssuer(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Service) ListIssuers(filter model.IssuerFilter, page, size int) ([]*model.Issuer, int, error) {
	all := s.store.ListIssuers()
	matched := make([]*model.Issuer, 0, len(all))
	for _, i := range all {
		if filter.Match(i) {
			matched = append(matched, i)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Name < matched[j].Name
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Issuer{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetIssuer(id string) (*model.Issuer, error) {
	return s.store.GetIssuer(id)
}

func (s *Service) UpdateIssuer(id string, input model.Issuer) (*model.Issuer, error) {
	i, err := s.store.GetIssuer(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		i.Name = input.Name
	}
	if input.Level != "" {
		i.Level = input.Level
	}
	i.ContactPhone = input.ContactPhone
	i.Address = input.Address
	if input.Status != "" {
		i.Status = input.Status
	}
	if err := i.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateIssuer(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Service) DeleteIssuer(id string) error {
	return s.store.DeleteIssuer(id)
}
