package store

import (
	"certmgmt/internal/model"
)

func (s *MemoryStore) CreateCertificate(c *model.Certificate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.certificates {
		if exist.Code == c.Code {
			return ErrConflict
		}
	}
	s.certificates[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCertificate(id string) (*model.Certificate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.certificates[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) GetCertificateByCode(code string) (*model.Certificate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.certificates {
		if c.Code == code {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListCertificates() []*model.Certificate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Certificate, 0, len(s.certificates))
	for _, c := range s.certificates {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateCertificate(c *model.Certificate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.certificates[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.certificates {
		if exist.ID != c.ID && exist.Code == c.Code {
			return ErrConflict
		}
	}
	s.certificates[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteCertificate(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.certificates[id]; !ok {
		return ErrNotFound
	}
	delete(s.certificates, id)
	return nil
}
