package service

import (
	"sort"
	"time"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateRenewal(input model.Renewal) (*model.Renewal, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCertificate(input.CertificateID); err != nil {
		return nil, model.NewValidationError("certificate_id", "证书不存在")
	}
	now := time.Now()
	r := &model.Renewal{
		ID:            idgen.Hex(),
		CertificateID: input.CertificateID,
		Applicant:     input.Applicant,
		Reason:        input.Reason,
		Status:        model.RenewalPending,
		ReviewComment: "",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.store.CreateRenewal(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) ListRenewals(filter model.RenewalFilter, page, size int) ([]*model.Renewal, int, error) {
	all := s.store.ListRenewals()
	matched := make([]*model.Renewal, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Renewal{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetRenewal(id string) (*model.Renewal, error) {
	return s.store.GetRenewal(id)
}

func (s *Service) UpdateRenewal(id string, input model.Renewal) (*model.Renewal, error) {
	r, err := s.store.GetRenewal(id)
	if err != nil {
		return nil, err
	}
	if input.Applicant != "" {
		r.Applicant = input.Applicant
	}
	r.Reason = input.Reason
	if input.Status != "" {
		if !model.RenewalCanTransition(r.Status, input.Status) {
			return nil, model.NewValidationError("status", "续期状态转换不合法")
		}
		r.Status = input.Status
	}
	r.ReviewComment = input.ReviewComment
	r.UpdatedAt = time.Now()
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateRenewal(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) ApproveRenewal(id string, comment string) (*model.Renewal, error) {
	r, err := s.store.GetRenewal(id)
	if err != nil {
		return nil, err
	}
	if !model.RenewalCanTransition(r.Status, model.RenewalApproved) {
		return nil, model.NewValidationError("status", "当前状态不能审批通过")
	}
	r.Status = model.RenewalApproved
	r.ReviewComment = comment
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateRenewal(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) RejectRenewal(id string, comment string) (*model.Renewal, error) {
	r, err := s.store.GetRenewal(id)
	if err != nil {
		return nil, err
	}
	if !model.RenewalCanTransition(r.Status, model.RenewalRejected) {
		return nil, model.NewValidationError("status", "当前状态不能驳回")
	}
	r.Status = model.RenewalRejected
	r.ReviewComment = comment
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateRenewal(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) CompleteRenewal(id string) (*model.Renewal, error) {
	r, err := s.store.GetRenewal(id)
	if err != nil {
		return nil, err
	}
	if !model.RenewalCanTransition(r.Status, model.RenewalCompleted) {
		return nil, model.NewValidationError("status", "当前状态不能完成")
	}
	r.Status = model.RenewalCompleted
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateRenewal(r); err != nil {
		return nil, err
	}
	cert, err := s.store.GetCertificate(r.CertificateID)
	if err == nil && cert.Status == model.CertificateExpired {
		cert.Status = model.CertificateActive
		cert.UpdatedAt = time.Now()
		_ = s.store.UpdateCertificate(cert)
	}
	return r, nil
}

func (s *Service) DeleteRenewal(id string) error {
	return s.store.DeleteRenewal(id)
}
