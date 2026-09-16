package store

import (
	"certmgmt/internal/model"
)

func (s *MemoryStore) CreateVerification(v *model.Verification) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.verifications[v.ID] = v
	return nil
}

func (s *MemoryStore) GetVerification(id string) (*model.Verification, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.verifications[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func (s *MemoryStore) ListVerifications() []*model.Verification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Verification, 0, len(s.verifications))
	for _, v := range s.verifications {
		list = append(list, v)
	}
	return list
}

func (s *MemoryStore) UpdateVerification(v *model.Verification) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.verifications[v.ID]; !ok {
		return ErrNotFound
	}
	s.verifications[v.ID] = v
	return nil
}

func (s *MemoryStore) DeleteVerification(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.verifications[id]; !ok {
		return ErrNotFound
	}
	delete(s.verifications, id)
	return nil
}
