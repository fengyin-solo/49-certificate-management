package handler

import (
	"net/http"
	"strconv"

	"certmgmt/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.getOverviewStats)
	mux.HandleFunc("GET /api/stats/by-category", s.getCategoryStats)
	mux.HandleFunc("GET /api/stats/by-status", s.getStatusStats)
	mux.HandleFunc("GET /api/stats/monthly-trend", s.getMonthlyTrend)
	mux.HandleFunc("GET /api/stats/top-issuers", s.getTopIssuers)
	mux.HandleFunc("GET /api/stats/expiring-soon", s.getExpiringSoon)
	mux.HandleFunc("GET /api/stats/export", s.exportSnapshot)
	mux.HandleFunc("GET /api/stats/training", s.getTrainingStats)
	mux.HandleFunc("GET /api/stats/reviews", s.getReviewStats)
	mux.HandleFunc("GET /api/stats/complaints", s.getComplaintStats)
}

func (s *Server) getOverviewStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetOverviewStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getCategoryStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetCategoryStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getStatusStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetStatusStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getMonthlyTrend(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetMonthlyTrend()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getTopIssuers(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	stats, err := s.svc.GetTopIssuers(limit)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getExpiringSoon(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	stats, err := s.svc.GetExpiringSoon(limit)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) exportSnapshot(w http.ResponseWriter, r *http.Request) {
	snapshot, err := s.svc.ExportSnapshot()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, snapshot)
}

func (s *Server) getTrainingStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetTrainingStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getReviewStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetReviewStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getComplaintStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetComplaintStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}
