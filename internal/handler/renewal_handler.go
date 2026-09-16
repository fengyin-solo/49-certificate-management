package handler

import (
	"net/http"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerRenewalRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/renewals", s.createRenewal)
	mux.HandleFunc("GET /api/renewals", s.listRenewals)
	mux.HandleFunc("GET /api/renewals/{id}", s.getRenewal)
	mux.HandleFunc("PUT /api/renewals/{id}", s.updateRenewal)
	mux.HandleFunc("DELETE /api/renewals/{id}", s.deleteRenewal)
	mux.HandleFunc("POST /api/renewals/{id}/approve", s.approveRenewal)
	mux.HandleFunc("POST /api/renewals/{id}/reject", s.rejectRenewal)
	mux.HandleFunc("POST /api/renewals/{id}/complete", s.completeRenewal)
}

type createRenewalRequest struct {
	CertificateID string `json:"certificate_id"`
	Applicant     string `json:"applicant"`
	Reason        string `json:"reason"`
}

func (s *Server) createRenewal(w http.ResponseWriter, r *http.Request) {
	var req createRenewalRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ren, err := s.svc.CreateRenewal(model.Renewal{
		CertificateID: req.CertificateID,
		Applicant:     req.Applicant,
		Reason:        req.Reason,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, ren)
}

func (s *Server) listRenewals(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RenewalFilter{
		CertificateID: r.URL.Query().Get("certificate_id"),
		Status:        r.URL.Query().Get("status"),
		Applicant:     r.URL.Query().Get("applicant"),
	}
	items, total, err := s.svc.ListRenewals(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRenewal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ren, err := s.svc.GetRenewal(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ren)
}

type updateRenewalRequest struct {
	Applicant     string `json:"applicant"`
	Reason        string `json:"reason"`
	Status        string `json:"status"`
	ReviewComment string `json:"review_comment"`
}

func (s *Server) updateRenewal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRenewalRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ren, err := s.svc.UpdateRenewal(id, model.Renewal{
		Applicant:     req.Applicant,
		Reason:        req.Reason,
		Status:        req.Status,
		ReviewComment: req.ReviewComment,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ren)
}

func (s *Server) deleteRenewal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRenewal(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type approveRenewalRequest struct {
	Comment string `json:"comment"`
}

func (s *Server) approveRenewal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req approveRenewalRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ren, err := s.svc.ApproveRenewal(id, req.Comment)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ren)
}

type rejectRenewalRequest struct {
	Comment string `json:"comment"`
}

func (s *Server) rejectRenewal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req rejectRenewalRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ren, err := s.svc.RejectRenewal(id, req.Comment)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ren)
}

func (s *Server) completeRenewal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ren, err := s.svc.CompleteRenewal(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ren)
}
