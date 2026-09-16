package handler

import (
	"net/http"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerIssuerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/issuers", s.createIssuer)
	mux.HandleFunc("GET /api/issuers", s.listIssuers)
	mux.HandleFunc("GET /api/issuers/{id}", s.getIssuer)
	mux.HandleFunc("PUT /api/issuers/{id}", s.updateIssuer)
	mux.HandleFunc("DELETE /api/issuers/{id}", s.deleteIssuer)
}

type createIssuerRequest struct {
	Name         string `json:"name"`
	Level        string `json:"level"`
	ContactPhone string `json:"contact_phone"`
	Address      string `json:"address"`
	Status       string `json:"status"`
}

func (s *Server) createIssuer(w http.ResponseWriter, r *http.Request) {
	var req createIssuerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	i, err := s.svc.CreateIssuer(model.Issuer{
		Name:         req.Name,
		Level:        req.Level,
		ContactPhone: req.ContactPhone,
		Address:      req.Address,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, i)
}

func (s *Server) listIssuers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.IssuerFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListIssuers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getIssuer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	i, err := s.svc.GetIssuer(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, i)
}

type updateIssuerRequest struct {
	Name         string `json:"name"`
	Level        string `json:"level"`
	ContactPhone string `json:"contact_phone"`
	Address      string `json:"address"`
	Status       string `json:"status"`
}

func (s *Server) updateIssuer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateIssuerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	i, err := s.svc.UpdateIssuer(id, model.Issuer{
		Name:         req.Name,
		Level:        req.Level,
		ContactPhone: req.ContactPhone,
		Address:      req.Address,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, i)
}

func (s *Server) deleteIssuer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteIssuer(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
