package store

import (
	"certmgmt/internal/model"
)

func (s *MemoryStore) CreateNotification(n *model.Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notifications[n.ID] = n
	return nil
}

func (s *MemoryStore) GetNotification(id string) (*model.Notification, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.notifications[id]
	if !ok {
		return nil, ErrNotFound
	}
	return n, nil
}

func (s *MemoryStore) ListNotifications() []*model.Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Notification, 0, len(s.notifications))
	for _, n := range s.notifications {
		list = append(list, n)
	}
	return list
}

func (s *MemoryStore) UpdateNotification(n *model.Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.notifications[n.ID]; !ok {
		return ErrNotFound
	}
	s.notifications[n.ID] = n
	return nil
}

func (s *MemoryStore) DeleteNotification(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.notifications[id]; !ok {
		return ErrNotFound
	}
	delete(s.notifications, id)
	return nil
}
