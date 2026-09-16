package model

import (
	"strings"
	"time"
)

type Attachment struct {
	ID            string    `json:"id"`
	CertificateID string    `json:"certificate_id"`
	FileName      string    `json:"file_name"`
	FileType      string    `json:"file_type"`
	Size          int64     `json:"size"`
	UploadedAt    time.Time `json:"uploaded_at"`
}

func (a *Attachment) Validate() error {
	a.CertificateID = strings.TrimSpace(a.CertificateID)
	a.FileName = strings.TrimSpace(a.FileName)
	a.FileType = strings.TrimSpace(a.FileType)
	if a.CertificateID == "" {
		return NewValidationError("certificate_id", "证书 ID 不能为空")
	}
	if a.FileName == "" {
		return NewValidationError("file_name", "文件名不能为空")
	}
	if a.FileType == "" {
		return NewValidationError("file_type", "文件类型不能为空")
	}
	if a.Size < 0 {
		return NewValidationError("size", "文件大小不能为负数")
	}
	return nil
}

type AttachmentFilter struct {
	CertificateID string
	FileType      string
}

func (f AttachmentFilter) Match(a *Attachment) bool {
	if f.CertificateID != "" && a.CertificateID != f.CertificateID {
		return false
	}
	if f.FileType != "" && a.FileType != f.FileType {
		return false
	}
	return true
}
