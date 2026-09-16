package store

import (
	"certmgmt/internal/model"
)

func (s *MemoryStore) CreateRevocation(r *model.Revocation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.revocations[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRevocation(id string) (*model.Revocation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.revocations[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListRevocations() []*model.Revocation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Revocation, 0, len(s.revocations))
	for _, r := range s.revocations {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRevocation(r *model.Revocation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.revocations[r.ID]; !ok {
		return ErrNotFound
	}
	s.revocations[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRevocation(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.revocations[id]; !ok {
		return ErrNotFound
	}
	delete(s.revocations, id)
	return nil
}
