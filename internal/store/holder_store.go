package store

import (
	"certmgmt/internal/model"
)

func (s *MemoryStore) CreateHolder(h *model.Holder) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.holders {
		if exist.IDNumber == h.IDNumber {
			return ErrConflict
		}
	}
	s.holders[h.ID] = h
	return nil
}

func (s *MemoryStore) GetHolder(id string) (*model.Holder, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h, ok := s.holders[id]
	if !ok {
		return nil, ErrNotFound
	}
	return h, nil
}

func (s *MemoryStore) ListHolders() []*model.Holder {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Holder, 0, len(s.holders))
	for _, h := range s.holders {
		list = append(list, h)
	}
	return list
}

func (s *MemoryStore) UpdateHolder(h *model.Holder) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.holders[h.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.holders {
		if exist.ID != h.ID && exist.IDNumber == h.IDNumber {
			return ErrConflict
		}
	}
	s.holders[h.ID] = h
	return nil
}

func (s *MemoryStore) DeleteHolder(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.holders[id]; !ok {
		return ErrNotFound
	}
	delete(s.holders, id)
	return nil
}
