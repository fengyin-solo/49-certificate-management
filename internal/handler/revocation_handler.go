package handler

import (
	"net/http"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerRevocationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/revocations", s.createRevocation)
	mux.HandleFunc("GET /api/revocations", s.listRevocations)
	mux.HandleFunc("GET /api/revocations/{id}", s.getRevocation)
	mux.HandleFunc("DELETE /api/revocations/{id}", s.deleteRevocation)
}

type createRevocationRequest struct {
	CertificateID string `json:"certificate_id"`
	Reason        string `json:"reason"`
	Operator      string `json:"operator"`
	Note          string `json:"note"`
}

func (s *Server) createRevocation(w http.ResponseWriter, r *http.Request) {
	var req createRevocationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rev, err := s.svc.CreateRevocation(model.Revocation{
		CertificateID: req.CertificateID,
		Reason:        req.Reason,
		Operator:      req.Operator,
		Note:          req.Note,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rev)
}

func (s *Server) listRevocations(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RevocationFilter{
		CertificateID: r.URL.Query().Get("certificate_id"),
		Operator:      r.URL.Query().Get("operator"),
	}
	items, total, err := s.svc.ListRevocations(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRevocation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rev, err := s.svc.GetRevocation(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rev)
}

func (s *Server) deleteRevocation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRevocation(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
