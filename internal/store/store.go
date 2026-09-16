// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"certmgmt/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	CreateCertificate(c *model.Certificate) error
	GetCertificate(id string) (*model.Certificate, error)
	GetCertificateByCode(code string) (*model.Certificate, error)
	ListCertificates() []*model.Certificate
	UpdateCertificate(c *model.Certificate) error
	DeleteCertificate(id string) error

	CreateCategory(c *model.Category) error
	GetCategory(id string) (*model.Category, error)
	ListCategories() []*model.Category
	UpdateCategory(c *model.Category) error
	DeleteCategory(id string) error

	CreateHolder(h *model.Holder) error
	GetHolder(id string) (*model.Holder, error)
	ListHolders() []*model.Holder
	UpdateHolder(h *model.Holder) error
	DeleteHolder(id string) error

	CreateIssuer(i *model.Issuer) error
	GetIssuer(id string) (*model.Issuer, error)
	ListIssuers() []*model.Issuer
	UpdateIssuer(i *model.Issuer) error
	DeleteIssuer(id string) error

	CreateVerification(v *model.Verification) error
	GetVerification(id string) (*model.Verification, error)
	ListVerifications() []*model.Verification
	UpdateVerification(v *model.Verification) error
	DeleteVerification(id string) error

	CreateReminder(r *model.Reminder) error
	GetReminder(id string) (*model.Reminder, error)
	ListReminders() []*model.Reminder
	UpdateReminder(r *model.Reminder) error
	DeleteReminder(id string) error

	CreateRenewal(r *model.Renewal) error
	GetRenewal(id string) (*model.Renewal, error)
	ListRenewals() []*model.Renewal
	UpdateRenewal(r *model.Renewal) error
	DeleteRenewal(id string) error

	CreateRevocation(r *model.Revocation) error
	GetRevocation(id string) (*model.Revocation, error)
	ListRevocations() []*model.Revocation
	UpdateRevocation(r *model.Revocation) error
	DeleteRevocation(id string) error

	CreateAttachment(a *model.Attachment) error
	GetAttachment(id string) (*model.Attachment, error)
	ListAttachments() []*model.Attachment
	UpdateAttachment(a *model.Attachment) error
	DeleteAttachment(id string) error

	CreateNotification(n *model.Notification) error
	GetNotification(id string) (*model.Notification, error)
	ListNotifications() []*model.Notification
	UpdateNotification(n *model.Notification) error
	DeleteNotification(id string) error

	CreateAuditLog(a *model.AuditLog) error
	GetAuditLog(id string) (*model.AuditLog, error)
	ListAuditLogs() []*model.AuditLog
	UpdateAuditLog(a *model.AuditLog) error
	DeleteAuditLog(id string) error

	CreateTrainingRecord(t *model.TrainingRecord) error
	GetTrainingRecord(id string) (*model.TrainingRecord, error)
	ListTrainingRecords() []*model.TrainingRecord
	UpdateTrainingRecord(t *model.TrainingRecord) error
	DeleteTrainingRecord(id string) error

	CreateReview(r *model.Review) error
	GetReview(id string) (*model.Review, error)
	ListReviews() []*model.Review
	UpdateReview(r *model.Review) error
	DeleteReview(id string) error

	CreateComplaint(c *model.Complaint) error
	GetComplaint(id string) (*model.Complaint, error)
	ListComplaints() []*model.Complaint
	UpdateComplaint(c *model.Complaint) error
	DeleteComplaint(id string) error
}
