package store

import (
	"sync"

	"certmgmt/internal/model"
)

type MemoryStore struct {
	mu              sync.RWMutex
	certificates    map[string]*model.Certificate
	categories      map[string]*model.Category
	holders         map[string]*model.Holder
	issuers         map[string]*model.Issuer
	verifications   map[string]*model.Verification
	reminders       map[string]*model.Reminder
	renewals        map[string]*model.Renewal
	revocations     map[string]*model.Revocation
	attachments     map[string]*model.Attachment
	notifications   map[string]*model.Notification
	auditLogs       map[string]*model.AuditLog
	trainingRecords map[string]*model.TrainingRecord
	reviews         map[string]*model.Review
	complaints      map[string]*model.Complaint
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		certificates:    make(map[string]*model.Certificate),
		categories:      make(map[string]*model.Category),
		holders:         make(map[string]*model.Holder),
		issuers:         make(map[string]*model.Issuer),
		verifications:   make(map[string]*model.Verification),
		reminders:       make(map[string]*model.Reminder),
		renewals:        make(map[string]*model.Renewal),
		revocations:     make(map[string]*model.Revocation),
		attachments:     make(map[string]*model.Attachment),
		notifications:   make(map[string]*model.Notification),
		auditLogs:       make(map[string]*model.AuditLog),
		trainingRecords: make(map[string]*model.TrainingRecord),
		reviews:         make(map[string]*model.Review),
		complaints:      make(map[string]*model.Complaint),
	}
}

var _ Store = (*MemoryStore)(nil)
