package service

import (
	"testing"
	"time"

	"certmgmt/internal/config"
	"certmgmt/internal/model"
	"certmgmt/internal/store"
	"certmgmt/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestService_CreateCertificate(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat A"})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})

	cert, err := svc.CreateCertificate(model.Certificate{
		Code: "C001", Name: "Cert A", CategoryID: cat.ID,
		HolderID: hold.ID, IssuerID: iss.ID,
		IssueDate: "2024-01-01", ExpireDate: "2025-01-01",
	})
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	if cert.ID == "" {
		t.Fatalf("certificate id empty")
	}
	if cert.Status != model.CertificateActive {
		t.Fatalf("expected active status")
	}

	_, err = svc.CreateCertificate(model.Certificate{Code: "C002", Name: "Bad", CategoryID: "bad", HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for bad category, got %v", err)
	}

	_, err = svc.CreateCertificate(model.Certificate{Code: "C001", Name: "Dup", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})
	if err != store.ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestService_CertificateStateMachine(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat"})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})
	cert, _ := svc.CreateCertificate(model.Certificate{Code: "C002", Name: "Cert", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})

	_, err := svc.UpdateCertificate(cert.ID, model.Certificate{Status: model.CertificateRevoked})
	if err != nil {
		t.Fatalf("active->revoked should succeed: %v", err)
	}

	_, err = svc.UpdateCertificate(cert.ID, model.Certificate{Status: model.CertificateActive})
	if !model.IsValidationError(err) {
		t.Fatalf("revoked->active should fail, got %v", err)
	}
}

func TestService_RevokeCertificateSideEffect(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat"})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})
	cert, _ := svc.CreateCertificate(model.Certificate{Code: "C003", Name: "Cert", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})

	cert, err := svc.RevokeCertificate(cert.ID, "Test reason", "Admin")
	if err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if cert.Status != model.CertificateRevoked {
		t.Fatalf("expected revoked")
	}
	revs, _, _ := svc.ListRevocations(model.RevocationFilter{CertificateID: cert.ID}, 1, 10)
	if len(revs) != 1 {
		t.Fatalf("expected 1 revocation record, got %d", len(revs))
	}
	logs, _, _ := svc.ListAuditLogs(model.AuditLogFilter{TargetID: cert.ID}, 1, 10)
	if len(logs) < 2 {
		t.Fatalf("expected audit logs for create+revoke")
	}
}

func TestService_BatchRevokeByHolder(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat"})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})
	_, _ = svc.CreateCertificate(model.Certificate{Code: "C004", Name: "Cert1", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})
	_, _ = svc.CreateCertificate(model.Certificate{Code: "C005", Name: "Cert2", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})

	count, err := svc.BatchRevokeByHolder(hold.ID, "Batch revoke", "Admin")
	if err != nil {
		t.Fatalf("batch revoke: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected revoke 2, got %d", count)
	}
}

func TestService_RenewalStateMachine(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat"})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})
	cert, _ := svc.CreateCertificate(model.Certificate{Code: "C006", Name: "Cert", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})
	ren, _ := svc.CreateRenewal(model.Renewal{CertificateID: cert.ID, Applicant: "Alice", Reason: "Expired"})

	_, err := svc.ApproveRenewal(ren.ID, "Approved")
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	ren, _ = svc.GetRenewal(ren.ID)
	if ren.Status != model.RenewalApproved {
		t.Fatalf("expected approved")
	}

	_, err = svc.RejectRenewal(ren.ID, "Should fail")
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for reject after approve, got %v", err)
	}

	cert, _ = svc.GetCertificate(cert.ID)
	cert.Status = model.CertificateExpired
	_ = svc.store.UpdateCertificate(cert)
	ren, err = svc.CompleteRenewal(ren.ID)
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if ren.Status != model.RenewalCompleted {
		t.Fatalf("expected completed")
	}
	cert, _ = svc.GetCertificate(cert.ID)
	if cert.Status != model.CertificateActive {
		t.Fatalf("expected certificate reactivated")
	}
}

func TestService_ReminderGeneration(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat", ValidYears: 3, RequireReview: true})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})
	cert, _ := svc.CreateCertificate(model.Certificate{Code: "C007", Name: "Cert", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})

	rems, _, _ := svc.ListReminders(model.ReminderFilter{CertificateID: cert.ID}, 1, 10)
	if len(rems) < 2 {
		t.Fatalf("expected at least 2 reminders, got %d", len(rems))
	}
}

func TestService_ListPagination(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat"})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})
	for i := 0; i < 5; i++ {
		_, _ = svc.CreateCertificate(model.Certificate{Code: "P" + string(rune('0'+i)), Name: "Cert", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})
	}
	items, total, err := svc.ListCertificates(model.CertificateFilter{}, 1, 2)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 5 {
		t.Fatalf("expected total 5, got %d", total)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	items, _, _ = svc.ListCertificates(model.CertificateFilter{}, 3, 2)
	if len(items) != 1 {
		t.Fatalf("expected 1 item on page 3, got %d", len(items))
	}
}

func TestService_Stats(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat"})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})
	_, _ = svc.CreateCertificate(model.Certificate{Code: "C008", Name: "Cert", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})

	ov, err := svc.GetOverviewStats()
	if err != nil {
		t.Fatalf("overview: %v", err)
	}
	if ov.TotalCertificates != 1 {
		t.Fatalf("expected 1 cert, got %d", ov.TotalCertificates)
	}

	cats, err := svc.GetCategoryStats()
	if err != nil {
		t.Fatalf("category stats: %v", err)
	}
	if len(cats) != 1 {
		t.Fatalf("expected 1 category group")
	}

	status, err := svc.GetStatusStats()
	if err != nil {
		t.Fatalf("status stats: %v", err)
	}
	if len(status) != 1 {
		t.Fatalf("expected 1 status group")
	}

	trend, err := svc.GetMonthlyTrend()
	if err != nil {
		t.Fatalf("monthly trend: %v", err)
	}
	if len(trend) != 1 {
		t.Fatalf("expected 1 month")
	}

	top, err := svc.GetTopIssuers(10)
	if err != nil {
		t.Fatalf("top issuers: %v", err)
	}
	if len(top) != 1 {
		t.Fatalf("expected 1 issuer")
	}

	exp, err := svc.GetExpiringSoon(10)
	if err != nil {
		t.Fatalf("expiring soon: %v", err)
	}
	_ = exp

	snap, err := svc.ExportSnapshot()
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if snap == nil {
		t.Fatalf("snapshot nil")
	}
}

func TestService_BatchImport(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat"})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})
	inputs := []model.Certificate{
		{Code: "B001", Name: "C1", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"},
		{Code: "B002", Name: "C2", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"},
	}
	certs, err := svc.BatchImportCertificates(inputs)
	if err != nil {
		t.Fatalf("batch import: %v", err)
	}
	if len(certs) != 2 {
		t.Fatalf("expected 2 certs, got %d", len(certs))
	}
}

func TestService_SendNotification(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat"})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})
	cert, _ := svc.CreateCertificate(model.Certificate{Code: "C009", Name: "Cert", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})
	rem, _ := svc.CreateReminder(model.Reminder{CertificateID: cert.ID, Type: model.ReminderBeforeExpire, TriggerAt: time.Now()})
	n, _ := svc.CreateNotification(model.Notification{ReminderID: rem.ID, Recipient: "Alice", Channel: model.ChannelEmail, Content: "Hello"})
	n2, err := svc.SendNotification(n.ID)
	if err != nil {
		t.Fatalf("send notification: %v", err)
	}
	if n2.Status != model.NotificationSent {
		t.Fatalf("expected sent")
	}
}

func TestService_TrainingRecord(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat"})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})
	cert, _ := svc.CreateCertificate(model.Certificate{Code: "C010", Name: "Cert", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})

	tr, err := svc.CreateTrainingRecord(model.TrainingRecord{CertificateID: cert.ID, Topic: "Safety", Trainer: "Bob", Hours: 8, TrainDate: "2024-06-01"})
	if err != nil {
		t.Fatalf("create training record: %v", err)
	}
	if tr.Status != model.TrainingPending {
		t.Fatalf("expected pending status")
	}

	_, err = svc.CreateTrainingRecord(model.TrainingRecord{CertificateID: "bad", Topic: "Safety", Trainer: "Bob", Hours: 8, TrainDate: "2024-06-01"})
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for bad certificate, got %v", err)
	}

	items, total, err := svc.ListTrainingRecords(model.TrainingRecordFilter{CertificateID: cert.ID}, 1, 10)
	if err != nil {
		t.Fatalf("list training records: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected 1 record, got %d", total)
	}

	tr, err = svc.UpdateTrainingRecord(tr.ID, model.TrainingRecord{Hours: 10, Status: model.TrainingCompleted})
	if err != nil {
		t.Fatalf("update training record: %v", err)
	}
	if tr.Hours != 10 {
		t.Fatalf("expected hours 10, got %d", tr.Hours)
	}

	report, err := svc.GetTrainingReport()
	if err != nil {
		t.Fatalf("training report: %v", err)
	}
	if len(report) != 1 {
		t.Fatalf("expected 1 report entry, got %d", len(report))
	}

	_, err = svc.ExportTrainingRecordsJSON()
	if err != nil {
		t.Fatalf("export training records: %v", err)
	}

	if err := svc.DeleteTrainingRecord(tr.ID); err != nil {
		t.Fatalf("delete training record: %v", err)
	}
}

func TestService_ReviewStateMachine(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat"})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})
	cert, _ := svc.CreateCertificate(model.Certificate{Code: "C011", Name: "Cert", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})

	rev, err := svc.CreateReview(model.Review{CertificateID: cert.ID, Cycle: 3, ReviewDate: "2024-06-01", Reviewer: "Alice"})
	if err != nil {
		t.Fatalf("create review: %v", err)
	}
	if rev.Status != model.ReviewPending {
		t.Fatalf("expected pending status")
	}

	_, err = svc.CreateReview(model.Review{CertificateID: "bad", Cycle: 3, ReviewDate: "2024-06-01", Reviewer: "Alice"})
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for bad certificate, got %v", err)
	}

	rev, err = svc.PassReview(rev.ID)
	if err != nil {
		t.Fatalf("pass review: %v", err)
	}
	if rev.Status != model.ReviewPassed {
		t.Fatalf("expected passed status")
	}

	_, err = svc.FailReview(rev.ID)
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for fail after pass, got %v", err)
	}

	items, total, err := svc.ListReviews(model.ReviewFilter{CertificateID: cert.ID}, 1, 10)
	if err != nil {
		t.Fatalf("list reviews: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected 1 review, got %d", total)
	}

	report, err := svc.GetReviewReport()
	if err != nil {
		t.Fatalf("review report: %v", err)
	}
	if len(report) != 1 {
		t.Fatalf("expected 1 report entry, got %d", len(report))
	}

	_, err = svc.ExportReviewsJSON()
	if err != nil {
		t.Fatalf("export reviews: %v", err)
	}

	if err := svc.DeleteReview(rev.ID); err != nil {
		t.Fatalf("delete review: %v", err)
	}
}

func TestService_ComplaintStateMachine(t *testing.T) {
	svc := newTestService()
	cat, _ := svc.CreateCategory(model.Category{Name: "Cat"})
	hold, _ := svc.CreateHolder(model.Holder{Name: "Alice", IDNumber: "ID001"})
	iss, _ := svc.CreateIssuer(model.Issuer{Name: "Org", Level: "L1"})
	cert, _ := svc.CreateCertificate(model.Certificate{Code: "C012", Name: "Cert", CategoryID: cat.ID, HolderID: hold.ID, IssuerID: iss.ID, IssueDate: "2024-01-01", ExpireDate: "2025-01-01"})

	comp, err := svc.CreateComplaint(model.Complaint{CertificateID: cert.ID, Complainant: "Alice", Type: "quality", Content: "bad"})
	if err != nil {
		t.Fatalf("create complaint: %v", err)
	}
	if comp.Status != model.ComplaintPending {
		t.Fatalf("expected pending status")
	}

	_, err = svc.CreateComplaint(model.Complaint{CertificateID: "bad", Complainant: "Alice", Type: "quality", Content: "bad"})
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for bad certificate, got %v", err)
	}

	comp, err = svc.InvestigateComplaint(comp.ID, "Bob")
	if err != nil {
		t.Fatalf("investigate complaint: %v", err)
	}
	if comp.Status != model.ComplaintInvestigating {
		t.Fatalf("expected investigating status")
	}

	_, err = svc.InvestigateComplaint(comp.ID, "Bob")
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for investigate again, got %v", err)
	}

	comp, err = svc.ResolveComplaint(comp.ID, "fixed")
	if err != nil {
		t.Fatalf("resolve complaint: %v", err)
	}
	if comp.Status != model.ComplaintResolved {
		t.Fatalf("expected resolved status")
	}

	_, err = svc.DismissComplaint(comp.ID, "dismiss")
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for dismiss after resolve, got %v", err)
	}

	items, total, err := svc.ListComplaints(model.ComplaintFilter{CertificateID: cert.ID}, 1, 10)
	if err != nil {
		t.Fatalf("list complaints: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected 1 complaint, got %d", total)
	}

	report, err := svc.GetComplaintReport()
	if err != nil {
		t.Fatalf("complaint report: %v", err)
	}
	if len(report) != 1 {
		t.Fatalf("expected 1 report entry, got %d", len(report))
	}

	_, err = svc.ExportComplaintsJSON()
	if err != nil {
		t.Fatalf("export complaints: %v", err)
	}

	if err := svc.DeleteComplaint(comp.ID); err != nil {
		t.Fatalf("delete complaint: %v", err)
	}
}
