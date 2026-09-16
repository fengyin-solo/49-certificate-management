package service

import (
	"sort"
	"time"

	"certmgmt/internal/model"
)

type OverviewStats struct {
	TotalCertificates   int `json:"total_certificates"`
	TotalCategories     int `json:"total_categories"`
	TotalHolders        int `json:"total_holders"`
	TotalIssuers        int `json:"total_issuers"`
	ActiveCertificates  int `json:"active_certificates"`
	ExpiredCertificates int `json:"expired_certificates"`
	RevokedCertificates int `json:"revoked_certificates"`
	PendingRenewals     int `json:"pending_renewals"`
}

func (s *Service) GetOverviewStats() (*OverviewStats, error) {
	stats := &OverviewStats{}
	for _, c := range s.store.ListCertificates() {
		stats.TotalCertificates++
		switch c.Status {
		case model.CertificateActive:
			stats.ActiveCertificates++
		case model.CertificateExpired:
			stats.ExpiredCertificates++
		case model.CertificateRevoked:
			stats.RevokedCertificates++
		}
	}
	stats.TotalCategories = len(s.store.ListCategories())
	stats.TotalHolders = len(s.store.ListHolders())
	stats.TotalIssuers = len(s.store.ListIssuers())
	for _, r := range s.store.ListRenewals() {
		if r.Status == model.RenewalPending {
			stats.PendingRenewals++
		}
	}
	return stats, nil
}

type CategoryGroup struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	Count        int    `json:"count"`
}

func (s *Service) GetCategoryStats() ([]CategoryGroup, error) {
	cats := s.store.ListCategories()
	catMap := make(map[string]string, len(cats))
	for _, c := range cats {
		catMap[c.ID] = c.Name
	}
	countMap := make(map[string]int)
	for _, c := range s.store.ListCertificates() {
		countMap[c.CategoryID]++
	}
	result := make([]CategoryGroup, 0, len(countMap))
	for id, count := range countMap {
		result = append(result, CategoryGroup{
			CategoryID:   id,
			CategoryName: catMap[id],
			Count:        count,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result, nil
}

type StatusGroup struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

func (s *Service) GetStatusStats() ([]StatusGroup, error) {
	countMap := make(map[string]int)
	for _, c := range s.store.ListCertificates() {
		countMap[c.Status]++
	}
	result := make([]StatusGroup, 0, len(countMap))
	for status, count := range countMap {
		result = append(result, StatusGroup{Status: status, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result, nil
}

type MonthlyTrend struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

func (s *Service) GetMonthlyTrend() ([]MonthlyTrend, error) {
	countMap := make(map[string]int)
	for _, c := range s.store.ListCertificates() {
		month := c.CreatedAt.Format("2006-01")
		countMap[month]++
	}
	result := make([]MonthlyTrend, 0, len(countMap))
	for month, count := range countMap {
		result = append(result, MonthlyTrend{Month: month, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Month < result[j].Month
	})
	return result, nil
}

type IssuerRanking struct {
	IssuerID   string `json:"issuer_id"`
	IssuerName string `json:"issuer_name"`
	Count      int    `json:"count"`
}

func (s *Service) GetTopIssuers(limit int) ([]IssuerRanking, error) {
	issuers := s.store.ListIssuers()
	issuerMap := make(map[string]string, len(issuers))
	for _, i := range issuers {
		issuerMap[i.ID] = i.Name
	}
	countMap := make(map[string]int)
	for _, c := range s.store.ListCertificates() {
		countMap[c.IssuerID]++
	}
	result := make([]IssuerRanking, 0, len(countMap))
	for id, count := range countMap {
		result = append(result, IssuerRanking{
			IssuerID:   id,
			IssuerName: issuerMap[id],
			Count:      count,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}
	return result, nil
}

type ExpiringSoon struct {
	CertificateID string `json:"certificate_id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	ExpireDate    string `json:"expire_date"`
	DaysLeft      int    `json:"days_left"`
}

func (s *Service) GetExpiringSoon(limit int) ([]ExpiringSoon, error) {
	now := time.Now()
	all := s.store.ListCertificates()
	result := make([]ExpiringSoon, 0)
	for _, c := range all {
		if c.Status == model.CertificateRevoked {
			continue
		}
		exp, err := time.Parse("2006-01-02", c.ExpireDate)
		if err != nil {
			continue
		}
		daysLeft := int(exp.Sub(now).Hours() / 24)
		if daysLeft <= 90 {
			result = append(result, ExpiringSoon{
				CertificateID: c.ID,
				Code:          c.Code,
				Name:          c.Name,
				ExpireDate:    c.ExpireDate,
				DaysLeft:      daysLeft,
			})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].DaysLeft < result[j].DaysLeft
	})
	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}
	return result, nil
}

type TrainingStats struct {
	CertificateID string `json:"certificate_id"`
	RecordCount   int    `json:"record_count"`
	TotalHours    int    `json:"total_hours"`
	PassCount     int    `json:"pass_count"`
}

func (s *Service) GetTrainingStats() ([]TrainingStats, error) {
	certMap := make(map[string]*TrainingStats)
	for _, t := range s.store.ListTrainingRecords() {
		st, ok := certMap[t.CertificateID]
		if !ok {
			st = &TrainingStats{CertificateID: t.CertificateID}
			certMap[t.CertificateID] = st
		}
		st.RecordCount++
		st.TotalHours += t.Hours
		if t.ExamResult == model.TrainingResultPass {
			st.PassCount++
		}
	}
	result := make([]TrainingStats, 0, len(certMap))
	for _, st := range certMap {
		result = append(result, *st)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].RecordCount > result[j].RecordCount
	})
	return result, nil
}

type ReviewStats struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

func (s *Service) GetReviewStats() ([]ReviewStats, error) {
	countMap := make(map[string]int)
	for _, r := range s.store.ListReviews() {
		countMap[r.Status]++
	}
	result := make([]ReviewStats, 0, len(countMap))
	for status, count := range countMap {
		result = append(result, ReviewStats{Status: status, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result, nil
}

type ComplaintStats struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	Count  int    `json:"count"`
}

func (s *Service) GetComplaintStats() ([]ComplaintStats, error) {
	type key struct{ Type, Status string }
	countMap := make(map[key]int)
	for _, c := range s.store.ListComplaints() {
		countMap[key{Type: c.Type, Status: c.Status}]++
	}
	result := make([]ComplaintStats, 0, len(countMap))
	for k, count := range countMap {
		result = append(result, ComplaintStats{Type: k.Type, Status: k.Status, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result, nil
}

func (s *Service) ExportSnapshot() (map[string]interface{}, error) {
	snapshot := map[string]interface{}{
		"certificates":     s.store.ListCertificates(),
		"categories":       s.store.ListCategories(),
		"holders":          s.store.ListHolders(),
		"issuers":          s.store.ListIssuers(),
		"verifications":    s.store.ListVerifications(),
		"reminders":        s.store.ListReminders(),
		"renewals":         s.store.ListRenewals(),
		"revocations":      s.store.ListRevocations(),
		"attachments":      s.store.ListAttachments(),
		"notifications":    s.store.ListNotifications(),
		"audit_logs":       s.store.ListAuditLogs(),
		"training_records": s.store.ListTrainingRecords(),
		"reviews":          s.store.ListReviews(),
		"complaints":       s.store.ListComplaints(),
		"exported_at":      time.Now().Format(time.RFC3339),
	}
	return snapshot, nil
}
