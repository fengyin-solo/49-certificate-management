package store

import (
	"certmgmt/internal/model"
)

func (s *MemoryStore) CreateRenewal(r *model.Renewal) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.renewals[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRenewal(id string) (*model.Renewal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.renewals[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListRenewals() []*model.Renewal {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Renewal, 0, len(s.renewals))
	for _, r := range s.renewals {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRenewal(r *model.Renewal) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.renewals[r.ID]; !ok {
		return ErrNotFound
	}
	s.renewals[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRenewal(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.renewals[id]; !ok {
		return ErrNotFound
	}
	delete(s.renewals, id)
	return nil
}
