package service

import (
	"sort"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateCategory(input model.Category) (*model.Category, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	cat := &model.Category{
		ID:            idgen.Hex(),
		Name:          input.Name,
		Description:   input.Description,
		ValidYears:    input.ValidYears,
		RequireReview: input.RequireReview,
		Status:        input.Status,
	}
	if err := s.store.CreateCategory(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *Service) ListCategories(filter model.CategoryFilter, page, size int) ([]*model.Category, int, error) {
	all := s.store.ListCategories()
	matched := make([]*model.Category, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Name < matched[j].Name
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Category{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetCategory(id string) (*model.Category, error) {
	return s.store.GetCategory(id)
}

func (s *Service) UpdateCategory(id string, input model.Category) (*model.Category, error) {
	cat, err := s.store.GetCategory(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		cat.Name = input.Name
	}
	cat.Description = input.Description
	if input.ValidYears >= 0 {
		cat.ValidYears = input.ValidYears
	}
	cat.RequireReview = input.RequireReview
	if input.Status != "" {
		cat.Status = input.Status
	}
	if err := cat.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCategory(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *Service) DeleteCategory(id string) error {
	return s.store.DeleteCategory(id)
}
