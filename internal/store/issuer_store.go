package store

import (
	"certmgmt/internal/model"
)

func (s *MemoryStore) CreateIssuer(i *model.Issuer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.issuers {
		if exist.Name == i.Name {
			return ErrConflict
		}
	}
	s.issuers[i.ID] = i
	return nil
}

func (s *MemoryStore) GetIssuer(id string) (*model.Issuer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.issuers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return i, nil
}

func (s *MemoryStore) ListIssuers() []*model.Issuer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Issuer, 0, len(s.issuers))
	for _, i := range s.issuers {
		list = append(list, i)
	}
	return list
}

func (s *MemoryStore) UpdateIssuer(i *model.Issuer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.issuers[i.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.issuers {
		if exist.ID != i.ID && exist.Name == i.Name {
			return ErrConflict
		}
	}
	s.issuers[i.ID] = i
	return nil
}

func (s *MemoryStore) DeleteIssuer(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.issuers[id]; !ok {
		return ErrNotFound
	}
	delete(s.issuers, id)
	return nil
}
