package handler

import (
	"net/http"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerHolderRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/holders", s.createHolder)
	mux.HandleFunc("GET /api/holders", s.listHolders)
	mux.HandleFunc("GET /api/holders/{id}", s.getHolder)
	mux.HandleFunc("PUT /api/holders/{id}", s.updateHolder)
	mux.HandleFunc("DELETE /api/holders/{id}", s.deleteHolder)
}

type createHolderRequest struct {
	Name     string `json:"name"`
	IDNumber string `json:"id_number"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Company  string `json:"company"`
	Status   string `json:"status"`
}

func (s *Server) createHolder(w http.ResponseWriter, r *http.Request) {
	var req createHolderRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	h, err := s.svc.CreateHolder(model.Holder{
		Name:     req.Name,
		IDNumber: req.IDNumber,
		Phone:    req.Phone,
		Email:    req.Email,
		Company:  req.Company,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, h)
}

func (s *Server) listHolders(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.HolderFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListHolders(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getHolder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	h, err := s.svc.GetHolder(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, h)
}

type updateHolderRequest struct {
	Name     string `json:"name"`
	IDNumber string `json:"id_number"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Company  string `json:"company"`
	Status   string `json:"status"`
}

func (s *Server) updateHolder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateHolderRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	h, err := s.svc.UpdateHolder(id, model.Holder{
		Name:     req.Name,
		IDNumber: req.IDNumber,
		Phone:    req.Phone,
		Email:    req.Email,
		Company:  req.Company,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, h)
}

func (s *Server) deleteHolder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteHolder(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
