package service

import (
	"sort"
	"time"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateComplaint(input model.Complaint) (*model.Complaint, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCertificate(input.CertificateID); err != nil {
		return nil, model.NewValidationError("certificate_id", "证书不存在")
	}
	now := time.Now()
	c := &model.Complaint{
		ID:            idgen.Hex(),
		CertificateID: input.CertificateID,
		Complainant:   input.Complainant,
		Type:          input.Type,
		Content:       input.Content,
		Status:        model.ComplaintPending,
		Handler:       input.Handler,
		HandleResult:  input.HandleResult,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.store.CreateComplaint(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListComplaints(filter model.ComplaintFilter, page, size int) ([]*model.Complaint, int, error) {
	all := s.store.ListComplaints()
	matched := make([]*model.Complaint, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Complaint{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetComplaint(id string) (*model.Complaint, error) {
	return s.store.GetComplaint(id)
}

func (s *Service) UpdateComplaint(id string, input model.Complaint) (*model.Complaint, error) {
	c, err := s.store.GetComplaint(id)
	if err != nil {
		return nil, err
	}
	if input.Type != "" {
		c.Type = input.Type
	}
	if input.Content != "" {
		c.Content = input.Content
	}
	if input.Handler != "" {
		c.Handler = input.Handler
	}
	if input.HandleResult != "" {
		c.HandleResult = input.HandleResult
	}
	if input.Status != "" {
		if !model.ComplaintCanTransition(c.Status, input.Status) {
			return nil, model.NewValidationError("status", "投诉状态转换不合法")
		}
		c.Status = input.Status
	}
	c.UpdatedAt = time.Now()
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateComplaint(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) InvestigateComplaint(id string, handler string) (*model.Complaint, error) {
	c, err := s.store.GetComplaint(id)
	if err != nil {
		return nil, err
	}
	if !model.ComplaintCanTransition(c.Status, model.ComplaintInvestigating) {
		return nil, model.NewValidationError("status", "当前状态不能进入调查中")
	}
	c.Status = model.ComplaintInvestigating
	c.Handler = handler
	c.UpdatedAt = time.Now()
	if err := s.store.UpdateComplaint(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ResolveComplaint(id string, result string) (*model.Complaint, error) {
	c, err := s.store.GetComplaint(id)
	if err != nil {
		return nil, err
	}
	if !model.ComplaintCanTransition(c.Status, model.ComplaintResolved) {
		return nil, model.NewValidationError("status", "当前状态不能标记为已解决")
	}
	c.Status = model.ComplaintResolved
	c.HandleResult = result
	c.UpdatedAt = time.Now()
	if err := s.store.UpdateComplaint(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DismissComplaint(id string, result string) (*model.Complaint, error) {
	c, err := s.store.GetComplaint(id)
	if err != nil {
		return nil, err
	}
	if !model.ComplaintCanTransition(c.Status, model.ComplaintDismissed) {
		return nil, model.NewValidationError("status", "当前状态不能标记为已驳回")
	}
	c.Status = model.ComplaintDismissed
	c.HandleResult = result
	c.UpdatedAt = time.Now()
	if err := s.store.UpdateComplaint(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteComplaint(id string) error {
	return s.store.DeleteComplaint(id)
}

func (s *Service) BatchImportComplaints(inputs []model.Complaint) ([]*model.Complaint, error) {
	results := make([]*model.Complaint, 0, len(inputs))
	for _, input := range inputs {
		c, err := s.CreateComplaint(input)
		if err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	return results, nil
}

func (s *Service) ExportComplaintsJSON() ([]byte, error) {
	all := s.store.ListComplaints()
	buf := make([]*model.Complaint, 0, len(all))
	for _, c := range all {
		buf = append(buf, c)
	}
	return jsonMarshalIndent(buf)
}

type ComplaintReport struct {
	Type           string `json:"type"`
	TotalCount     int    `json:"total_count"`
	ResolvedCount  int    `json:"resolved_count"`
	DismissedCount int    `json:"dismissed_count"`
	PendingCount   int    `json:"pending_count"`
}

func (s *Service) GetComplaintReport() ([]ComplaintReport, error) {
	typeMap := make(map[string]*ComplaintReport)
	for _, c := range s.store.ListComplaints() {
		rep, ok := typeMap[c.Type]
		if !ok {
			rep = &ComplaintReport{Type: c.Type}
			typeMap[c.Type] = rep
		}
		rep.TotalCount++
		switch c.Status {
		case model.ComplaintResolved:
			rep.ResolvedCount++
		case model.ComplaintDismissed:
			rep.DismissedCount++
		case model.ComplaintPending, model.ComplaintInvestigating:
			rep.PendingCount++
		}
	}
	result := make([]ComplaintReport, 0, len(typeMap))
	for _, r := range typeMap {
		result = append(result, *r)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalCount > result[j].TotalCount
	})
	return result, nil
}
