package handler

import (
	"net/http"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerNotificationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/notifications", s.createNotification)
	mux.HandleFunc("GET /api/notifications", s.listNotifications)
	mux.HandleFunc("GET /api/notifications/{id}", s.getNotification)
	mux.HandleFunc("POST /api/notifications/{id}/send", s.sendNotification)
	mux.HandleFunc("DELETE /api/notifications/{id}", s.deleteNotification)
}

type createNotificationRequest struct {
	ReminderID string `json:"reminder_id"`
	Recipient  string `json:"recipient"`
	Channel    string `json:"channel"`
	Status     string `json:"status"`
	Content    string `json:"content"`
}

func (s *Server) createNotification(w http.ResponseWriter, r *http.Request) {
	var req createNotificationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.CreateNotification(model.Notification{
		ReminderID: req.ReminderID,
		Recipient:  req.Recipient,
		Channel:    req.Channel,
		Status:     req.Status,
		Content:    req.Content,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, n)
}

func (s *Server) listNotifications(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.NotificationFilter{
		ReminderID: r.URL.Query().Get("reminder_id"),
		Channel:    r.URL.Query().Get("channel"),
		Status:     r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListNotifications(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.GetNotification(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) sendNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.SendNotification(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) deleteNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteNotification(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
