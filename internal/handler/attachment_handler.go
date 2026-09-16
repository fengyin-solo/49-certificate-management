package handler

import (
	"net/http"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerAttachmentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/attachments", s.createAttachment)
	mux.HandleFunc("GET /api/attachments", s.listAttachments)
	mux.HandleFunc("GET /api/attachments/{id}", s.getAttachment)
	mux.HandleFunc("DELETE /api/attachments/{id}", s.deleteAttachment)
}

type createAttachmentRequest struct {
	CertificateID string `json:"certificate_id"`
	FileName      string `json:"file_name"`
	FileType      string `json:"file_type"`
	Size          int64  `json:"size"`
}

func (s *Server) createAttachment(w http.ResponseWriter, r *http.Request) {
	var req createAttachmentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAttachment(model.Attachment{
		CertificateID: req.CertificateID,
		FileName:      req.FileName,
		FileType:      req.FileType,
		Size:          req.Size,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAttachments(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AttachmentFilter{
		CertificateID: r.URL.Query().Get("certificate_id"),
		FileType:      r.URL.Query().Get("file_type"),
	}
	items, total, err := s.svc.ListAttachments(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAttachment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.GetAttachment(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAttachment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAttachment(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
