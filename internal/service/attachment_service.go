package service

import (
	"sort"
	"time"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateAttachment(input model.Attachment) (*model.Attachment, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCertificate(input.CertificateID); err != nil {
		return nil, model.NewValidationError("certificate_id", "证书不存在")
	}
	a := &model.Attachment{
		ID:            idgen.Hex(),
		CertificateID: input.CertificateID,
		FileName:      input.FileName,
		FileType:      input.FileType,
		Size:          input.Size,
		UploadedAt:    time.Now(),
	}
	if err := s.store.CreateAttachment(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) ListAttachments(filter model.AttachmentFilter, page, size int) ([]*model.Attachment, int, error) {
	all := s.store.ListAttachments()
	matched := make([]*model.Attachment, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].UploadedAt.After(matched[j].UploadedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Attachment{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetAttachment(id string) (*model.Attachment, error) {
	return s.store.GetAttachment(id)
}

func (s *Service) DeleteAttachment(id string) error {
	return s.store.DeleteAttachment(id)
}
