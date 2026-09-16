package handler

import (
	"net/http"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerCertificateRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/certificates", s.createCertificate)
	mux.HandleFunc("GET /api/certificates", s.listCertificates)
	mux.HandleFunc("GET /api/certificates/{id}", s.getCertificate)
	mux.HandleFunc("PUT /api/certificates/{id}", s.updateCertificate)
	mux.HandleFunc("DELETE /api/certificates/{id}", s.deleteCertificate)
	mux.HandleFunc("POST /api/certificates/{id}/revoke", s.revokeCertificate)
	mux.HandleFunc("POST /api/certificates/batch-import", s.batchImportCertificates)
	mux.HandleFunc("POST /api/certificates/batch-revoke", s.batchRevokeByHolder)
}

type createCertificateRequest struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	CategoryID string `json:"category_id"`
	HolderID   string `json:"holder_id"`
	IssuerID   string `json:"issuer_id"`
	IssueDate  string `json:"issue_date"`
	ExpireDate string `json:"expire_date"`
	Status     string `json:"status"`
	Remarks    string `json:"remarks"`
}

func (s *Server) createCertificate(w http.ResponseWriter, r *http.Request) {
	var req createCertificateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	cert, err := s.svc.CreateCertificate(model.Certificate{
		Code:       req.Code,
		Name:       req.Name,
		CategoryID: req.CategoryID,
		HolderID:   req.HolderID,
		IssuerID:   req.IssuerID,
		IssueDate:  req.IssueDate,
		ExpireDate: req.ExpireDate,
		Status:     req.Status,
		Remarks:    req.Remarks,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, cert)
}

func (s *Server) listCertificates(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CertificateFilter{
		CategoryID: r.URL.Query().Get("category_id"),
		HolderID:   r.URL.Query().Get("holder_id"),
		IssuerID:   r.URL.Query().Get("issuer_id"),
		Status:     r.URL.Query().Get("status"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListCertificates(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCertificate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cert, err := s.svc.GetCertificate(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, cert)
}

type updateCertificateRequest struct {
	Name       string `json:"name"`
	CategoryID string `json:"category_id"`
	HolderID   string `json:"holder_id"`
	IssuerID   string `json:"issuer_id"`
	IssueDate  string `json:"issue_date"`
	ExpireDate string `json:"expire_date"`
	Status     string `json:"status"`
	Remarks    string `json:"remarks"`
}

func (s *Server) updateCertificate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateCertificateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	cert, err := s.svc.UpdateCertificate(id, model.Certificate{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		HolderID:   req.HolderID,
		IssuerID:   req.IssuerID,
		IssueDate:  req.IssueDate,
		ExpireDate: req.ExpireDate,
		Status:     req.Status,
		Remarks:    req.Remarks,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, cert)
}

func (s *Server) deleteCertificate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCertificate(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type revokeCertificateRequest struct {
	Reason   string `json:"reason"`
	Operator string `json:"operator"`
}

func (s *Server) revokeCertificate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req revokeCertificateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	cert, err := s.svc.RevokeCertificate(id, req.Reason, req.Operator)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, cert)
}

type batchImportRequest struct {
	Items []createCertificateRequest `json:"items"`
}

func (s *Server) batchImportCertificates(w http.ResponseWriter, r *http.Request) {
	var req batchImportRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.Certificate, len(req.Items))
	for i, item := range req.Items {
		inputs[i] = model.Certificate{
			Code:       item.Code,
			Name:       item.Name,
			CategoryID: item.CategoryID,
			HolderID:   item.HolderID,
			IssuerID:   item.IssuerID,
			IssueDate:  item.IssueDate,
			ExpireDate: item.ExpireDate,
			Status:     item.Status,
			Remarks:    item.Remarks,
		}
	}
	certs, err := s.svc.BatchImportCertificates(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, certs)
}

type batchRevokeRequest struct {
	HolderID string `json:"holder_id"`
	Reason   string `json:"reason"`
	Operator string `json:"operator"`
}

func (s *Server) batchRevokeByHolder(w http.ResponseWriter, r *http.Request) {
	var req batchRevokeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	count, err := s.svc.BatchRevokeByHolder(req.HolderID, req.Reason, req.Operator)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"revoked_count": count})
}
