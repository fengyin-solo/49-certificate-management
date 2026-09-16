package store

import (
	"certmgmt/internal/model"
)

func (s *MemoryStore) CreateReminder(r *model.Reminder) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reminders[r.ID] = r
	return nil
}

func (s *MemoryStore) GetReminder(id string) (*model.Reminder, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.reminders[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListReminders() []*model.Reminder {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Reminder, 0, len(s.reminders))
	for _, r := range s.reminders {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateReminder(r *model.Reminder) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.reminders[r.ID]; !ok {
		return ErrNotFound
	}
	s.reminders[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteReminder(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.reminders[id]; !ok {
		return ErrNotFound
	}
	delete(s.reminders, id)
	return nil
}
