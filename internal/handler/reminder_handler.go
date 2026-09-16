package handler

import (
	"net/http"
	"time"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerReminderRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/reminders", s.createReminder)
	mux.HandleFunc("GET /api/reminders", s.listReminders)
	mux.HandleFunc("GET /api/reminders/{id}", s.getReminder)
	mux.HandleFunc("PUT /api/reminders/{id}", s.updateReminder)
	mux.HandleFunc("DELETE /api/reminders/{id}", s.deleteReminder)
}

type createReminderRequest struct {
	CertificateID string    `json:"certificate_id"`
	Type          string    `json:"type"`
	AdvanceDays   int       `json:"advance_days"`
	TriggerAt     time.Time `json:"trigger_at"`
	Status        string    `json:"status"`
	Message       string    `json:"message"`
}

func (s *Server) createReminder(w http.ResponseWriter, r *http.Request) {
	var req createReminderRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rem, err := s.svc.CreateReminder(model.Reminder{
		CertificateID: req.CertificateID,
		Type:          req.Type,
		AdvanceDays:   req.AdvanceDays,
		TriggerAt:     req.TriggerAt,
		Status:        req.Status,
		Message:       req.Message,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rem)
}

func (s *Server) listReminders(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ReminderFilter{
		CertificateID: r.URL.Query().Get("certificate_id"),
		Type:          r.URL.Query().Get("type"),
		Status:        r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListReminders(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getReminder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rem, err := s.svc.GetReminder(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rem)
}

type updateReminderRequest struct {
	Type        string    `json:"type"`
	AdvanceDays int       `json:"advance_days"`
	TriggerAt   time.Time `json:"trigger_at"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
}

func (s *Server) updateReminder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateReminderRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rem, err := s.svc.UpdateReminder(id, model.Reminder{
		Type:        req.Type,
		AdvanceDays: req.AdvanceDays,
		TriggerAt:   req.TriggerAt,
		Status:      req.Status,
		Message:     req.Message,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rem)
}

func (s *Server) deleteReminder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteReminder(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
