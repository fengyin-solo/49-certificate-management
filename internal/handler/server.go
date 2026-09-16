// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"certmgmt/internal/config"
	"certmgmt/internal/model"
	"certmgmt/internal/service"
	"certmgmt/internal/store"
	"certmgmt/pkg/httpx"
	"certmgmt/pkg/logger"
)

type Server struct {
	svc *service.Service
	log *logger.Logger
	cfg *config.Config
}

func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, log: log, cfg: cfg}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerCertificateRoutes(mux)
	s.registerCategoryRoutes(mux)
	s.registerHolderRoutes(mux)
	s.registerIssuerRoutes(mux)
	s.registerVerificationRoutes(mux)
	s.registerReminderRoutes(mux)
	s.registerRenewalRoutes(mux)
	s.registerRevocationRoutes(mux)
	s.registerAttachmentRoutes(mux)
	s.registerNotificationRoutes(mux)
	s.registerAuditLogRoutes(mux)
	s.registerStatsRoutes(mux)
	s.registerTrainingRecordRoutes(mux)
	s.registerReviewRoutes(mux)
	s.registerComplaintRoutes(mux)
	mux.Handle("GET /", http.FileServer(http.Dir("web")))
	return s.withMiddleware(mux)
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.APIKey == "" {
			next.ServeHTTP(w, r)
			return
		}
		key := r.Header.Get("X-API-Key")
		if key == "" {
			key = r.URL.Query().Get("api_key")
		}
		if key != s.cfg.APIKey {
			httpx.Unauthorized(w, "API Key 无效")
			return
		}
		next.ServeHTTP(w, r)
	})
}

type rateLimitEntry struct {
	count  int
	window time.Time
}

var rateLimitMap = make(map[string]*rateLimitEntry)
var rateLimitMu sync.Mutex

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if ip == "" {
			ip = r.RemoteAddr
		}
		now := time.Now()
		rateLimitMu.Lock()
		entry, ok := rateLimitMap[ip]
		if !ok || now.Sub(entry.window) > time.Minute {
			entry = &rateLimitEntry{count: 0, window: now}
			rateLimitMap[ip] = entry
		}
		entry.count++
		count := entry.count
		rateLimitMu.Unlock()
		if count > 120 {
			httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁，请稍后再试")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}
