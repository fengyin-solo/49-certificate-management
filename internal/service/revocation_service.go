package service

import (
	"sort"
	"time"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateRevocation(input model.Revocation) (*model.Revocation, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCertificate(input.CertificateID); err != nil {
		return nil, model.NewValidationError("certificate_id", "证书不存在")
	}
	r := &model.Revocation{
		ID:            idgen.Hex(),
		CertificateID: input.CertificateID,
		Reason:        input.Reason,
		Operator:      input.Operator,
		RevokedAt:     time.Now(),
		Note:          input.Note,
	}
	if err := s.store.CreateRevocation(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) ListRevocations(filter model.RevocationFilter, page, size int) ([]*model.Revocation, int, error) {
	all := s.store.ListRevocations()
	matched := make([]*model.Revocation, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].RevokedAt.After(matched[j].RevokedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Revocation{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetRevocation(id string) (*model.Revocation, error) {
	return s.store.GetRevocation(id)
}

func (s *Service) DeleteRevocation(id string) error {
	return s.store.DeleteRevocation(id)
}
