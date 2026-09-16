package store

import (
	"certmgmt/internal/model"
)

func (s *MemoryStore) CreateComplaint(c *model.Complaint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.complaints[c.ID] = c
	return nil
}

func (s *MemoryStore) GetComplaint(id string) (*model.Complaint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.complaints[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) ListComplaints() []*model.Complaint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Complaint, 0, len(s.complaints))
	for _, c := range s.complaints {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateComplaint(c *model.Complaint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.complaints[c.ID]; !ok {
		return ErrNotFound
	}
	s.complaints[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteComplaint(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.complaints[id]; !ok {
		return ErrNotFound
	}
	delete(s.complaints, id)
	return nil
}
