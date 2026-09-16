package handler

import (
	"net/http"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerReviewRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/reviews", s.createReview)
	mux.HandleFunc("GET /api/reviews", s.listReviews)
	mux.HandleFunc("GET /api/reviews/{id}", s.getReview)
	mux.HandleFunc("PUT /api/reviews/{id}", s.updateReview)
	mux.HandleFunc("DELETE /api/reviews/{id}", s.deleteReview)
	mux.HandleFunc("POST /api/reviews/{id}/pass", s.passReview)
	mux.HandleFunc("POST /api/reviews/{id}/fail", s.failReview)
	mux.HandleFunc("POST /api/reviews/batch-import", s.batchImportReviews)
	mux.HandleFunc("GET /api/reviews/export", s.exportReviews)
	mux.HandleFunc("GET /api/reviews/report", s.getReviewReport)
}

type createReviewRequest struct {
	CertificateID  string `json:"certificate_id"`
	Cycle          int    `json:"cycle"`
	ReviewDate     string `json:"review_date"`
	Result         string `json:"result"`
	NextReviewDate string `json:"next_review_date"`
	Reviewer       string `json:"reviewer"`
}

func (s *Server) createReview(w http.ResponseWriter, r *http.Request) {
	var req createReviewRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rev, err := s.svc.CreateReview(model.Review{
		CertificateID:  req.CertificateID,
		Cycle:          req.Cycle,
		ReviewDate:     req.ReviewDate,
		Result:         req.Result,
		NextReviewDate: req.NextReviewDate,
		Reviewer:       req.Reviewer,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rev)
}

func (s *Server) listReviews(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ReviewFilter{
		CertificateID: r.URL.Query().Get("certificate_id"),
		Status:        r.URL.Query().Get("status"),
		Reviewer:      r.URL.Query().Get("reviewer"),
	}
	items, total, err := s.svc.ListReviews(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getReview(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rev, err := s.svc.GetReview(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rev)
}

type updateReviewRequest struct {
	Cycle          int    `json:"cycle"`
	ReviewDate     string `json:"review_date"`
	Result         string `json:"result"`
	NextReviewDate string `json:"next_review_date"`
	Reviewer       string `json:"reviewer"`
	Status         string `json:"status"`
}

func (s *Server) updateReview(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateReviewRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rev, err := s.svc.UpdateReview(id, model.Review{
		Cycle:          req.Cycle,
		ReviewDate:     req.ReviewDate,
		Result:         req.Result,
		NextReviewDate: req.NextReviewDate,
		Reviewer:       req.Reviewer,
		Status:         req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rev)
}

func (s *Server) deleteReview(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteReview(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) passReview(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rev, err := s.svc.PassReview(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rev)
}

func (s *Server) failReview(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rev, err := s.svc.FailReview(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rev)
}

func (s *Server) batchImportReviews(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Items []createReviewRequest `json:"items"`
	}
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.Review, len(req.Items))
	for i, item := range req.Items {
		inputs[i] = model.Review{
			CertificateID:  item.CertificateID,
			Cycle:          item.Cycle,
			ReviewDate:     item.ReviewDate,
			Result:         item.Result,
			NextReviewDate: item.NextReviewDate,
			Reviewer:       item.Reviewer,
		}
	}
	items, err := s.svc.BatchImportReviews(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, items)
}

func (s *Server) exportReviews(w http.ResponseWriter, r *http.Request) {
	data, err := s.svc.ExportReviewsJSON()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=reviews.json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) getReviewReport(w http.ResponseWriter, r *http.Request) {
	report, err := s.svc.GetReviewReport()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, report)
}
