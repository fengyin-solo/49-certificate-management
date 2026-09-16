package handler

import (
	"net/http"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerAuditLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/audit-logs", s.createAuditLog)
	mux.HandleFunc("GET /api/audit-logs", s.listAuditLogs)
	mux.HandleFunc("GET /api/audit-logs/{id}", s.getAuditLog)
	mux.HandleFunc("DELETE /api/audit-logs/{id}", s.deleteAuditLog)
}

type createAuditLogRequest struct {
	Operator   string `json:"operator"`
	Action     string `json:"action"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	Detail     string `json:"detail"`
}

func (s *Server) createAuditLog(w http.ResponseWriter, r *http.Request) {
	var req createAuditLogRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAuditLog(model.AuditLog{
		Operator:   req.Operator,
		Action:     req.Action,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Detail:     req.Detail,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAuditLogs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AuditLogFilter{
		Operator:   r.URL.Query().Get("operator"),
		Action:     r.URL.Query().Get("action"),
		TargetType: r.URL.Query().Get("target_type"),
		TargetID:   r.URL.Query().Get("target_id"),
	}
	items, total, err := s.svc.ListAuditLogs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAuditLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.GetAuditLog(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAuditLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAuditLog(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
