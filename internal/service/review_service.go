package service

import (
	"sort"
	"time"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateReview(input model.Review) (*model.Review, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCertificate(input.CertificateID); err != nil {
		return nil, model.NewValidationError("certificate_id", "证书不存在")
	}
	now := time.Now()
	r := &model.Review{
		ID:             idgen.Hex(),
		CertificateID:  input.CertificateID,
		Cycle:          input.Cycle,
		ReviewDate:     input.ReviewDate,
		Result:         input.Result,
		NextReviewDate: input.NextReviewDate,
		Reviewer:       input.Reviewer,
		Status:         model.ReviewPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.store.CreateReview(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) ListReviews(filter model.ReviewFilter, page, size int) ([]*model.Review, int, error) {
	all := s.store.ListReviews()
	matched := make([]*model.Review, 0, len(all))
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
		return []*model.Review{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetReview(id string) (*model.Review, error) {
	return s.store.GetReview(id)
}

func (s *Service) UpdateReview(id string, input model.Review) (*model.Review, error) {
	r, err := s.store.GetReview(id)
	if err != nil {
		return nil, err
	}
	if input.Cycle > 0 {
		r.Cycle = input.Cycle
	}
	if input.ReviewDate != "" {
		r.ReviewDate = input.ReviewDate
	}
	if input.NextReviewDate != "" {
		r.NextReviewDate = input.NextReviewDate
	}
	if input.Reviewer != "" {
		r.Reviewer = input.Reviewer
	}
	if input.Status != "" {
		if !model.ReviewCanTransition(r.Status, input.Status) {
			return nil, model.NewValidationError("status", "复审状态转换不合法")
		}
		r.Status = input.Status
		r.Result = input.Result
	}
	r.UpdatedAt = time.Now()
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateReview(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) PassReview(id string) (*model.Review, error) {
	r, err := s.store.GetReview(id)
	if err != nil {
		return nil, err
	}
	if !model.ReviewCanTransition(r.Status, model.ReviewPassed) {
		return nil, model.NewValidationError("status", "当前状态不能标记为通过")
	}
	r.Status = model.ReviewPassed
	r.Result = model.ReviewPassed
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateReview(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) FailReview(id string) (*model.Review, error) {
	r, err := s.store.GetReview(id)
	if err != nil {
		return nil, err
	}
	if !model.ReviewCanTransition(r.Status, model.ReviewFailed) {
		return nil, model.NewValidationError("status", "当前状态不能标记为不通过")
	}
	r.Status = model.ReviewFailed
	r.Result = model.ReviewFailed
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateReview(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteReview(id string) error {
	return s.store.DeleteReview(id)
}

func (s *Service) BatchImportReviews(inputs []model.Review) ([]*model.Review, error) {
	results := make([]*model.Review, 0, len(inputs))
	for _, input := range inputs {
		rev, err := s.CreateReview(input)
		if err != nil {
			return nil, err
		}
		results = append(results, rev)
	}
	return results, nil
}

func (s *Service) ExportReviewsJSON() ([]byte, error) {
	all := s.store.ListReviews()
	buf := make([]*model.Review, 0, len(all))
	for _, r := range all {
		buf = append(buf, r)
	}
	return jsonMarshalIndent(buf)
}

type ReviewReport struct {
	CertificateID string `json:"certificate_id"`
	TotalCount    int    `json:"total_count"`
	PassedCount   int    `json:"passed_count"`
	FailedCount   int    `json:"failed_count"`
	PendingCount  int    `json:"pending_count"`
}

func (s *Service) GetReviewReport() ([]ReviewReport, error) {
	certMap := make(map[string]*ReviewReport)
	for _, r := range s.store.ListReviews() {
		rep, ok := certMap[r.CertificateID]
		if !ok {
			rep = &ReviewReport{CertificateID: r.CertificateID}
			certMap[r.CertificateID] = rep
		}
		rep.TotalCount++
		switch r.Status {
		case model.ReviewPassed:
			rep.PassedCount++
		case model.ReviewFailed:
			rep.FailedCount++
		case model.ReviewPending:
			rep.PendingCount++
		}
	}
	result := make([]ReviewReport, 0, len(certMap))
	for _, r := range certMap {
		result = append(result, *r)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalCount > result[j].TotalCount
	})
	return result, nil
}
