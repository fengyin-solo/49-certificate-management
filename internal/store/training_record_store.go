package store

import (
	"certmgmt/internal/model"
)

func (s *MemoryStore) CreateTrainingRecord(t *model.TrainingRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.trainingRecords[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTrainingRecord(id string) (*model.TrainingRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.trainingRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListTrainingRecords() []*model.TrainingRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TrainingRecord, 0, len(s.trainingRecords))
	for _, t := range s.trainingRecords {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTrainingRecord(t *model.TrainingRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.trainingRecords[t.ID]; !ok {
		return ErrNotFound
	}
	s.trainingRecords[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTrainingRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.trainingRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.trainingRecords, id)
	return nil
}
