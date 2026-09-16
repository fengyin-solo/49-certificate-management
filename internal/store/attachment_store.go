package store

import (
	"certmgmt/internal/model"
)

func (s *MemoryStore) CreateAttachment(a *model.Attachment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attachments[a.ID] = a
	return nil
}

func (s *MemoryStore) GetAttachment(id string) (*model.Attachment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.attachments[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) ListAttachments() []*model.Attachment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Attachment, 0, len(s.attachments))
	for _, a := range s.attachments {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) UpdateAttachment(a *model.Attachment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.attachments[a.ID]; !ok {
		return ErrNotFound
	}
	s.attachments[a.ID] = a
	return nil
}

func (s *MemoryStore) DeleteAttachment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.attachments[id]; !ok {
		return ErrNotFound
	}
	delete(s.attachments, id)
	return nil
}
