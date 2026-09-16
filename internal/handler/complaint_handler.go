package handler

import (
	"net/http"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerComplaintRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/complaints", s.createComplaint)
	mux.HandleFunc("GET /api/complaints", s.listComplaints)
	mux.HandleFunc("GET /api/complaints/{id}", s.getComplaint)
	mux.HandleFunc("PUT /api/complaints/{id}", s.updateComplaint)
	mux.HandleFunc("DELETE /api/complaints/{id}", s.deleteComplaint)
	mux.HandleFunc("POST /api/complaints/{id}/investigate", s.investigateComplaint)
	mux.HandleFunc("POST /api/complaints/{id}/resolve", s.resolveComplaint)
	mux.HandleFunc("POST /api/complaints/{id}/dismiss", s.dismissComplaint)
	mux.HandleFunc("POST /api/complaints/batch-import", s.batchImportComplaints)
	mux.HandleFunc("GET /api/complaints/export", s.exportComplaints)
	mux.HandleFunc("GET /api/complaints/report", s.getComplaintReport)
}

type createComplaintRequest struct {
	CertificateID string `json:"certificate_id"`
	Complainant   string `json:"complainant"`
	Type          string `json:"type"`
	Content       string `json:"content"`
	Handler       string `json:"handler"`
	HandleResult  string `json:"handle_result"`
}

func (s *Server) createComplaint(w http.ResponseWriter, r *http.Request) {
	var req createComplaintRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateComplaint(model.Complaint{
		CertificateID: req.CertificateID,
		Complainant:   req.Complainant,
		Type:          req.Type,
		Content:       req.Content,
		Handler:       req.Handler,
		HandleResult:  req.HandleResult,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listComplaints(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ComplaintFilter{
		CertificateID: r.URL.Query().Get("certificate_id"),
		Status:        r.URL.Query().Get("status"),
		Type:          r.URL.Query().Get("type"),
		Complainant:   r.URL.Query().Get("complainant"),
	}
	items, total, err := s.svc.ListComplaints(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getComplaint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetComplaint(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type updateComplaintRequest struct {
	Type         string `json:"type"`
	Content      string `json:"content"`
	Status       string `json:"status"`
	Handler      string `json:"handler"`
	HandleResult string `json:"handle_result"`
}

func (s *Server) updateComplaint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateComplaintRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateComplaint(id, model.Complaint{
		Type:         req.Type,
		Content:      req.Content,
		Status:       req.Status,
		Handler:      req.Handler,
		HandleResult: req.HandleResult,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteComplaint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteComplaint(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type investigateComplaintRequest struct {
	Handler string `json:"handler"`
}

func (s *Server) investigateComplaint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req investigateComplaintRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.InvestigateComplaint(id, req.Handler)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type resolveComplaintRequest struct {
	Result string `json:"result"`
}

func (s *Server) resolveComplaint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req resolveComplaintRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.ResolveComplaint(id, req.Result)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type dismissComplaintRequest struct {
	Result string `json:"result"`
}

func (s *Server) dismissComplaint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dismissComplaintRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.DismissComplaint(id, req.Result)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) batchImportComplaints(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Items []createComplaintRequest `json:"items"`
	}
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.Complaint, len(req.Items))
	for i, item := range req.Items {
		inputs[i] = model.Complaint{
			CertificateID: item.CertificateID,
			Complainant:   item.Complainant,
			Type:          item.Type,
			Content:       item.Content,
			Handler:       item.Handler,
			HandleResult:  item.HandleResult,
		}
	}
	items, err := s.svc.BatchImportComplaints(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, items)
}

func (s *Server) exportComplaints(w http.ResponseWriter, r *http.Request) {
	data, err := s.svc.ExportComplaintsJSON()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=complaints.json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) getComplaintReport(w http.ResponseWriter, r *http.Request) {
	report, err := s.svc.GetComplaintReport()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, report)
}
