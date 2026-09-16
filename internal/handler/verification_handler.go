package handler

import (
	"net/http"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerVerificationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/verifications", s.createVerification)
	mux.HandleFunc("GET /api/verifications", s.listVerifications)
	mux.HandleFunc("GET /api/verifications/{id}", s.getVerification)
	mux.HandleFunc("PUT /api/verifications/{id}", s.updateVerification)
	mux.HandleFunc("DELETE /api/verifications/{id}", s.deleteVerification)
}

type createVerificationRequest struct {
	CertificateID string `json:"certificate_id"`
	Method        string `json:"method"`
	Result        string `json:"result"`
	Verifier      string `json:"verifier"`
	Note          string `json:"note"`
}

func (s *Server) createVerification(w http.ResponseWriter, r *http.Request) {
	var req createVerificationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.CreateVerification(model.Verification{
		CertificateID: req.CertificateID,
		Method:        req.Method,
		Result:        req.Result,
		Verifier:      req.Verifier,
		Note:          req.Note,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, v)
}

func (s *Server) listVerifications(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.VerificationFilter{
		CertificateID: r.URL.Query().Get("certificate_id"),
		Method:        r.URL.Query().Get("method"),
		Result:        r.URL.Query().Get("result"),
	}
	items, total, err := s.svc.ListVerifications(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getVerification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	v, err := s.svc.GetVerification(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

type updateVerificationRequest struct {
	Method   string `json:"method"`
	Result   string `json:"result"`
	Verifier string `json:"verifier"`
	Note     string `json:"note"`
}

func (s *Server) updateVerification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateVerificationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.UpdateVerification(id, model.Verification{
		Method:   req.Method,
		Result:   req.Result,
		Verifier: req.Verifier,
		Note:     req.Note,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) deleteVerification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteVerification(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
