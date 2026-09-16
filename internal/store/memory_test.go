package store

import (
	"testing"

	"certmgmt/internal/model"
)

func TestMemoryStore_Certificate(t *testing.T) {
	s := NewMemoryStore()
	c := &model.Certificate{ID: "c1", Code: "C001", Name: "Test", CategoryID: "cat1", HolderID: "h1", IssuerID: "i1", IssueDate: "2024-01-01", ExpireDate: "2025-01-01", Status: model.CertificateActive}
	if err := s.CreateCertificate(c); err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	if err := s.CreateCertificate(c); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetCertificate("c1")
	if err != nil {
		t.Fatalf("get certificate: %v", err)
	}
	if got.Code != "C001" {
		t.Fatalf("code mismatch")
	}
	if _, err := s.GetCertificate("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListCertificates()) != 1 {
		t.Fatalf("list count mismatch")
	}
	c.Name = "Updated"
	if err := s.UpdateCertificate(c); err != nil {
		t.Fatalf("update certificate: %v", err)
	}
	if err := s.UpdateCertificate(&model.Certificate{ID: "x"}); err != ErrNotFound {
		t.Fatalf("expected not found on update, got %v", err)
	}
	if err := s.DeleteCertificate("c1"); err != nil {
		t.Fatalf("delete certificate: %v", err)
	}
	if err := s.DeleteCertificate("c1"); err != ErrNotFound {
		t.Fatalf("expected not found on delete, got %v", err)
	}
}

func TestMemoryStore_Category(t *testing.T) {
	s := NewMemoryStore()
	c := &model.Category{ID: "cat1", Name: "Category A"}
	if err := s.CreateCategory(c); err != nil {
		t.Fatalf("create category: %v", err)
	}
	if err := s.CreateCategory(c); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	if _, err := s.GetCategory("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	c.Name = "Category B"
	if err := s.UpdateCategory(c); err != nil {
		t.Fatalf("update category: %v", err)
	}
	if err := s.DeleteCategory("cat1"); err != nil {
		t.Fatalf("delete category: %v", err)
	}
}

func TestMemoryStore_Holder(t *testing.T) {
	s := NewMemoryStore()
	h := &model.Holder{ID: "h1", Name: "Alice", IDNumber: "ID001"}
	if err := s.CreateHolder(h); err != nil {
		t.Fatalf("create holder: %v", err)
	}
	if err := s.CreateHolder(h); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	if _, err := s.GetHolder("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteHolder("h1"); err != nil {
		t.Fatalf("delete holder: %v", err)
	}
}

func TestMemoryStore_Issuer(t *testing.T) {
	s := NewMemoryStore()
	i := &model.Issuer{ID: "i1", Name: "Org A", Level: "L1"}
	if err := s.CreateIssuer(i); err != nil {
		t.Fatalf("create issuer: %v", err)
	}
	if err := s.CreateIssuer(i); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	if _, err := s.GetIssuer("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteIssuer("i1"); err != nil {
		t.Fatalf("delete issuer: %v", err)
	}
}

func TestMemoryStore_Verification(t *testing.T) {
	s := NewMemoryStore()
	v := &model.Verification{ID: "v1", CertificateID: "c1", Method: model.VerificationOnline, Result: model.VerifyValid, Verifier: "Bob"}
	if err := s.CreateVerification(v); err != nil {
		t.Fatalf("create verification: %v", err)
	}
	if _, err := s.GetVerification("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteVerification("v1"); err != nil {
		t.Fatalf("delete verification: %v", err)
	}
}

func TestMemoryStore_Reminder(t *testing.T) {
	s := NewMemoryStore()
	r := &model.Reminder{ID: "r1", CertificateID: "c1", Type: model.ReminderBeforeExpire}
	if err := s.CreateReminder(r); err != nil {
		t.Fatalf("create reminder: %v", err)
	}
	if _, err := s.GetReminder("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteReminder("r1"); err != nil {
		t.Fatalf("delete reminder: %v", err)
	}
}

func TestMemoryStore_Renewal(t *testing.T) {
	s := NewMemoryStore()
	ren := &model.Renewal{ID: "ren1", CertificateID: "c1", Applicant: "Alice"}
	if err := s.CreateRenewal(ren); err != nil {
		t.Fatalf("create renewal: %v", err)
	}
	if _, err := s.GetRenewal("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteRenewal("ren1"); err != nil {
		t.Fatalf("delete renewal: %v", err)
	}
}

func TestMemoryStore_Revocation(t *testing.T) {
	s := NewMemoryStore()
	rev := &model.Revocation{ID: "rev1", CertificateID: "c1", Reason: "Expired", Operator: "Admin"}
	if err := s.CreateRevocation(rev); err != nil {
		t.Fatalf("create revocation: %v", err)
	}
	if _, err := s.GetRevocation("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteRevocation("rev1"); err != nil {
		t.Fatalf("delete revocation: %v", err)
	}
}

func TestMemoryStore_Attachment(t *testing.T) {
	s := NewMemoryStore()
	a := &model.Attachment{ID: "a1", CertificateID: "c1", FileName: "doc.pdf", FileType: "pdf", Size: 1024}
	if err := s.CreateAttachment(a); err != nil {
		t.Fatalf("create attachment: %v", err)
	}
	if _, err := s.GetAttachment("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteAttachment("a1"); err != nil {
		t.Fatalf("delete attachment: %v", err)
	}
}

func TestMemoryStore_Notification(t *testing.T) {
	s := NewMemoryStore()
	n := &model.Notification{ID: "n1", ReminderID: "r1", Recipient: "Alice", Channel: model.ChannelEmail}
	if err := s.CreateNotification(n); err != nil {
		t.Fatalf("create notification: %v", err)
	}
	if _, err := s.GetNotification("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteNotification("n1"); err != nil {
		t.Fatalf("delete notification: %v", err)
	}
}

func TestMemoryStore_AuditLog(t *testing.T) {
	s := NewMemoryStore()
	a := &model.AuditLog{ID: "al1", Operator: "Admin", Action: "create", TargetType: "certificate", TargetID: "c1"}
	if err := s.CreateAuditLog(a); err != nil {
		t.Fatalf("create audit log: %v", err)
	}
	if _, err := s.GetAuditLog("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteAuditLog("al1"); err != nil {
		t.Fatalf("delete audit log: %v", err)
	}
}

func TestMemoryStore_TrainingRecord(t *testing.T) {
	s := NewMemoryStore()
	tr := &model.TrainingRecord{ID: "tr1", CertificateID: "c1", Topic: "Safety", Trainer: "Bob", Hours: 8, TrainDate: "2024-01-01"}
	if err := s.CreateTrainingRecord(tr); err != nil {
		t.Fatalf("create training record: %v", err)
	}
	got, err := s.GetTrainingRecord("tr1")
	if err != nil {
		t.Fatalf("get training record: %v", err)
	}
	if got.Topic != "Safety" {
		t.Fatalf("topic mismatch")
	}
	if _, err := s.GetTrainingRecord("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	tr.Topic = "Updated"
	if err := s.UpdateTrainingRecord(tr); err != nil {
		t.Fatalf("update training record: %v", err)
	}
	if err := s.UpdateTrainingRecord(&model.TrainingRecord{ID: "x"}); err != ErrNotFound {
		t.Fatalf("expected not found on update, got %v", err)
	}
	if err := s.DeleteTrainingRecord("tr1"); err != nil {
		t.Fatalf("delete training record: %v", err)
	}
	if err := s.DeleteTrainingRecord("tr1"); err != ErrNotFound {
		t.Fatalf("expected not found on delete, got %v", err)
	}
}

func TestMemoryStore_Review(t *testing.T) {
	s := NewMemoryStore()
	rev := &model.Review{ID: "rev1", CertificateID: "c1", Cycle: 3, ReviewDate: "2024-06-01", Reviewer: "Alice"}
	if err := s.CreateReview(rev); err != nil {
		t.Fatalf("create review: %v", err)
	}
	got, err := s.GetReview("rev1")
	if err != nil {
		t.Fatalf("get review: %v", err)
	}
	if got.Cycle != 3 {
		t.Fatalf("cycle mismatch")
	}
	if _, err := s.GetReview("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	rev.Reviewer = "Updated"
	if err := s.UpdateReview(rev); err != nil {
		t.Fatalf("update review: %v", err)
	}
	if err := s.UpdateReview(&model.Review{ID: "x"}); err != ErrNotFound {
		t.Fatalf("expected not found on update, got %v", err)
	}
	if err := s.DeleteReview("rev1"); err != nil {
		t.Fatalf("delete review: %v", err)
	}
	if err := s.DeleteReview("rev1"); err != ErrNotFound {
		t.Fatalf("expected not found on delete, got %v", err)
	}
}

func TestMemoryStore_Complaint(t *testing.T) {
	s := NewMemoryStore()
	c := &model.Complaint{ID: "comp1", CertificateID: "c1", Complainant: "Alice", Type: "quality", Content: "bad"}
	if err := s.CreateComplaint(c); err != nil {
		t.Fatalf("create complaint: %v", err)
	}
	got, err := s.GetComplaint("comp1")
	if err != nil {
		t.Fatalf("get complaint: %v", err)
	}
	if got.Type != "quality" {
		t.Fatalf("type mismatch")
	}
	if _, err := s.GetComplaint("notfound"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	c.Content = "updated"
	if err := s.UpdateComplaint(c); err != nil {
		t.Fatalf("update complaint: %v", err)
	}
	if err := s.UpdateComplaint(&model.Complaint{ID: "x"}); err != ErrNotFound {
		t.Fatalf("expected not found on update, got %v", err)
	}
	if err := s.DeleteComplaint("comp1"); err != nil {
		t.Fatalf("delete complaint: %v", err)
	}
	if err := s.DeleteComplaint("comp1"); err != ErrNotFound {
		t.Fatalf("expected not found on delete, got %v", err)
	}
}
