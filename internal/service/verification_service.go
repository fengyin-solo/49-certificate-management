package service

import (
	"sort"
	"time"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateVerification(input model.Verification) (*model.Verification, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCertificate(input.CertificateID); err != nil {
		return nil, model.NewValidationError("certificate_id", "证书不存在")
	}
	v := &model.Verification{
		ID:            idgen.Hex(),
		CertificateID: input.CertificateID,
		Method:        input.Method,
		Result:        input.Result,
		Verifier:      input.Verifier,
		VerifiedAt:    time.Now(),
		Note:          input.Note,
	}
	if err := s.store.CreateVerification(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Service) ListVerifications(filter model.VerificationFilter, page, size int) ([]*model.Verification, int, error) {
	all := s.store.ListVerifications()
	matched := make([]*model.Verification, 0, len(all))
	for _, v := range all {
		if filter.Match(v) {
			matched = append(matched, v)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].VerifiedAt.After(matched[j].VerifiedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Verification{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetVerification(id string) (*model.Verification, error) {
	return s.store.GetVerification(id)
}

func (s *Service) UpdateVerification(id string, input model.Verification) (*model.Verification, error) {
	v, err := s.store.GetVerification(id)
	if err != nil {
		return nil, err
	}
	if input.Method != "" {
		v.Method = input.Method
	}
	if input.Result != "" {
		v.Result = input.Result
	}
	if input.Verifier != "" {
		v.Verifier = input.Verifier
	}
	v.Note = input.Note
	if err := v.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateVerification(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Service) DeleteVerification(id string) error {
	return s.store.DeleteVerification(id)
}
