package model

import (
	"strings"
	"time"
)

const (
	TrainingResultPass = "pass"
	TrainingResultFail = "fail"
)

const (
	TrainingCompleted = "completed"
	TrainingPending   = "pending"
)

type TrainingRecord struct {
	ID            string    `json:"id"`
	CertificateID string    `json:"certificate_id"`
	Topic         string    `json:"topic"`
	Trainer       string    `json:"trainer"`
	Hours         int       `json:"hours"`
	TrainDate     string    `json:"train_date"`
	ExamResult    string    `json:"exam_result"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

func (t *TrainingRecord) Validate() error {
	t.CertificateID = strings.TrimSpace(t.CertificateID)
	t.Topic = strings.TrimSpace(t.Topic)
	t.Trainer = strings.TrimSpace(t.Trainer)
	t.TrainDate = strings.TrimSpace(t.TrainDate)
	if t.CertificateID == "" {
		return NewValidationError("certificate_id", "证书 ID 不能为空")
	}
	if t.Topic == "" {
		return NewValidationError("topic", "培训主题不能为空")
	}
	if t.Trainer == "" {
		return NewValidationError("trainer", "培训师不能为空")
	}
	if t.Hours <= 0 {
		return NewValidationError("hours", "培训时长必须大于 0")
	}
	if t.TrainDate == "" {
		return NewValidationError("train_date", "培训日期不能为空")
	}
	if t.ExamResult == "" {
		t.ExamResult = TrainingResultFail
	}
	if t.ExamResult != TrainingResultPass && t.ExamResult != TrainingResultFail {
		return NewValidationError("exam_result", "考试结果不合法")
	}
	if t.Status == "" {
		t.Status = TrainingPending
	}
	if t.Status != TrainingCompleted && t.Status != TrainingPending {
		return NewValidationError("status", "培训状态不合法")
	}
	return nil
}

type TrainingRecordFilter struct {
	CertificateID string
	Status        string
	Topic         string
}

func (f TrainingRecordFilter) Match(t *TrainingRecord) bool {
	if f.CertificateID != "" && t.CertificateID != f.CertificateID {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Topic != "" {
		k := strings.ToLower(strings.TrimSpace(f.Topic))
		if k != "" && !strings.Contains(strings.ToLower(t.Topic), k) {
			return false
		}
	}
	return true
}
