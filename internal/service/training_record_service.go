package service

import (
	"sort"
	"time"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateTrainingRecord(input model.TrainingRecord) (*model.TrainingRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCertificate(input.CertificateID); err != nil {
		return nil, model.NewValidationError("certificate_id", "证书不存在")
	}
	now := time.Now()
	t := &model.TrainingRecord{
		ID:            idgen.Hex(),
		CertificateID: input.CertificateID,
		Topic:         input.Topic,
		Trainer:       input.Trainer,
		Hours:         input.Hours,
		TrainDate:     input.TrainDate,
		ExamResult:    input.ExamResult,
		Status:        input.Status,
		CreatedAt:     now,
	}
	if err := s.store.CreateTrainingRecord(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ListTrainingRecords(filter model.TrainingRecordFilter, page, size int) ([]*model.TrainingRecord, int, error) {
	all := s.store.ListTrainingRecords()
	matched := make([]*model.TrainingRecord, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.TrainingRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetTrainingRecord(id string) (*model.TrainingRecord, error) {
	return s.store.GetTrainingRecord(id)
}

func (s *Service) UpdateTrainingRecord(id string, input model.TrainingRecord) (*model.TrainingRecord, error) {
	t, err := s.store.GetTrainingRecord(id)
	if err != nil {
		return nil, err
	}
	if input.Topic != "" {
		t.Topic = input.Topic
	}
	if input.Trainer != "" {
		t.Trainer = input.Trainer
	}
	if input.Hours > 0 {
		t.Hours = input.Hours
	}
	if input.TrainDate != "" {
		t.TrainDate = input.TrainDate
	}
	if input.ExamResult != "" {
		t.ExamResult = input.ExamResult
	}
	if input.Status != "" {
		t.Status = input.Status
	}
	if err := t.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateTrainingRecord(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteTrainingRecord(id string) error {
	return s.store.DeleteTrainingRecord(id)
}

func (s *Service) BatchImportTrainingRecords(inputs []model.TrainingRecord) ([]*model.TrainingRecord, error) {
	results := make([]*model.TrainingRecord, 0, len(inputs))
	for _, input := range inputs {
		tr, err := s.CreateTrainingRecord(input)
		if err != nil {
			return nil, err
		}
		results = append(results, tr)
	}
	return results, nil
}

func (s *Service) ExportTrainingRecordsJSON() ([]byte, error) {
	all := s.store.ListTrainingRecords()
	buf := make([]*model.TrainingRecord, 0, len(all))
	for _, t := range all {
		buf = append(buf, t)
	}
	return jsonMarshalIndent(buf)
}

type TrainingReport struct {
	CertificateID string `json:"certificate_id"`
	RecordCount   int    `json:"record_count"`
	TotalHours    int    `json:"total_hours"`
	PassRate      string `json:"pass_rate"`
}

func (s *Service) GetTrainingReport() ([]TrainingReport, error) {
	certMap := make(map[string]*TrainingReport)
	for _, t := range s.store.ListTrainingRecords() {
		r, ok := certMap[t.CertificateID]
		if !ok {
			r = &TrainingReport{CertificateID: t.CertificateID}
			certMap[t.CertificateID] = r
		}
		r.RecordCount++
		r.TotalHours += t.Hours
	}
	result := make([]TrainingReport, 0, len(certMap))
	for _, r := range certMap {
		result = append(result, *r)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].RecordCount > result[j].RecordCount
	})
	return result, nil
}
