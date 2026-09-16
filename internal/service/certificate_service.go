package service

import (
	"sort"
	"time"

	"certmgmt/internal/model"
	"certmgmt/pkg/idgen"
)

func (s *Service) CreateCertificate(input model.Certificate) (*model.Certificate, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCategory(input.CategoryID); err != nil {
		return nil, model.NewValidationError("category_id", "分类不存在")
	}
	if _, err := s.store.GetHolder(input.HolderID); err != nil {
		return nil, model.NewValidationError("holder_id", "持有人不存在")
	}
	if _, err := s.store.GetIssuer(input.IssuerID); err != nil {
		return nil, model.NewValidationError("issuer_id", "发证机构不存在")
	}
	now := time.Now()
	cert := &model.Certificate{
		ID:         idgen.Hex(),
		Code:       input.Code,
		Name:       input.Name,
		CategoryID: input.CategoryID,
		HolderID:   input.HolderID,
		IssuerID:   input.IssuerID,
		IssueDate:  input.IssueDate,
		ExpireDate: input.ExpireDate,
		Status:     input.Status,
		Remarks:    input.Remarks,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateCertificate(cert); err != nil {
		return nil, err
	}
	s.generateReminderForCertificate(cert)
	s.createAuditLog("create_certificate", "certificate", cert.ID, "创建证书")
	return cert, nil
}

func (s *Service) ListCertificates(filter model.CertificateFilter, page, size int) ([]*model.Certificate, int, error) {
	all := s.store.ListCertificates()
	matched := make([]*model.Certificate, 0, len(all))
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
		return []*model.Certificate{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetCertificate(id string) (*model.Certificate, error) {
	return s.store.GetCertificate(id)
}

func (s *Service) UpdateCertificate(id string, input model.Certificate) (*model.Certificate, error) {
	cert, err := s.store.GetCertificate(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		cert.Name = input.Name
	}
	if input.CategoryID != "" {
		cert.CategoryID = input.CategoryID
	}
	if input.HolderID != "" {
		cert.HolderID = input.HolderID
	}
	if input.IssuerID != "" {
		cert.IssuerID = input.IssuerID
	}
	if input.IssueDate != "" {
		cert.IssueDate = input.IssueDate
	}
	if input.ExpireDate != "" {
		cert.ExpireDate = input.ExpireDate
	}
	if input.Status != "" {
		if !model.CertificateCanTransition(cert.Status, input.Status) {
			return nil, model.NewValidationError("status", "状态转换不合法")
		}
		cert.Status = input.Status
	}
	cert.Remarks = input.Remarks
	cert.UpdatedAt = time.Now()
	if err := cert.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCertificate(cert); err != nil {
		return nil, err
	}
	s.createAuditLog("update_certificate", "certificate", cert.ID, "更新证书")
	return cert, nil
}

func (s *Service) DeleteCertificate(id string) error {
	if _, err := s.store.GetCertificate(id); err != nil {
		return err
	}
	s.createAuditLog("delete_certificate", "certificate", id, "删除证书")
	return s.store.DeleteCertificate(id)
}

func (s *Service) RevokeCertificate(id string, reason string, operator string) (*model.Certificate, error) {
	cert, err := s.store.GetCertificate(id)
	if err != nil {
		return nil, err
	}
	if !model.CertificateCanTransition(cert.Status, model.CertificateRevoked) {
		return nil, model.NewValidationError("status", "当前状态不能吊销")
	}
	cert.Status = model.CertificateRevoked
	cert.UpdatedAt = time.Now()
	if err := s.store.UpdateCertificate(cert); err != nil {
		return nil, err
	}
	rev := &model.Revocation{
		ID:            idgen.Hex(),
		CertificateID: id,
		Reason:        reason,
		Operator:      operator,
		RevokedAt:     time.Now(),
		Note:          "",
	}
	if err := s.store.CreateRevocation(rev); err != nil {
		return nil, err
	}
	s.createAuditLog("revoke_certificate", "certificate", id, "吊销证书: "+reason)
	return cert, nil
}

func (s *Service) BatchRevokeByHolder(holderID string, reason string, operator string) (int, error) {
	all := s.store.ListCertificates()
	count := 0
	for _, c := range all {
		if c.HolderID == holderID && c.Status != model.CertificateRevoked {
			if model.CertificateCanTransition(c.Status, model.CertificateRevoked) {
				c.Status = model.CertificateRevoked
				c.UpdatedAt = time.Now()
				if err := s.store.UpdateCertificate(c); err == nil {
					rev := &model.Revocation{
						ID:            idgen.Hex(),
						CertificateID: c.ID,
						Reason:        reason,
						Operator:      operator,
						RevokedAt:     time.Now(),
					}
					_ = s.store.CreateRevocation(rev)
					count++
				}
			}
		}
	}
	if count > 0 {
		s.createAuditLog("batch_revoke", "holder", holderID, "批量吊销 "+string(rune('0'+count))+" 张证书")
	}
	return count, nil
}

func (s *Service) BatchImportCertificates(inputs []model.Certificate) ([]*model.Certificate, error) {
	results := make([]*model.Certificate, 0, len(inputs))
	for _, input := range inputs {
		cert, err := s.CreateCertificate(input)
		if err != nil {
			return nil, err
		}
		results = append(results, cert)
	}
	return results, nil
}

func (s *Service) generateReminderForCertificate(cert *model.Certificate) {
	cat, err := s.store.GetCategory(cert.CategoryID)
	if err != nil {
		return
	}
	issueTime, _ := time.Parse("2006-01-02", cert.IssueDate)
	expireTime, _ := time.Parse("2006-01-02", cert.ExpireDate)
	if cat.ValidYears > 0 && cat.RequireReview {
		triggerAt := issueTime.AddDate(cat.ValidYears, 0, 0)
		rem := &model.Reminder{
			ID:            idgen.Hex(),
			CertificateID: cert.ID,
			Type:          model.ReminderReviewDue,
			AdvanceDays:   0,
			TriggerAt:     triggerAt,
			Status:        model.ReminderPending,
			Message:       "证书复审提醒",
		}
		_ = s.store.CreateReminder(rem)
	}
	rem := &model.Reminder{
		ID:            idgen.Hex(),
		CertificateID: cert.ID,
		Type:          model.ReminderBeforeExpire,
		AdvanceDays:   30,
		TriggerAt:     expireTime.AddDate(0, 0, -30),
		Status:        model.ReminderPending,
		Message:       "证书即将到期提醒",
	}
	_ = s.store.CreateReminder(rem)
}

func (s *Service) createAuditLog(action, targetType, targetID, detail string) {
	log := &model.AuditLog{
		ID:         idgen.Hex(),
		Operator:   "system",
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Detail:     detail,
		CreatedAt:  time.Now(),
	}
	_ = s.store.CreateAuditLog(log)
}
