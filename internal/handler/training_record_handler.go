package handler

import (
	"net/http"

	"certmgmt/internal/model"
	"certmgmt/pkg/httpx"
)

func (s *Server) registerTrainingRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/training-records", s.createTrainingRecord)
	mux.HandleFunc("GET /api/training-records", s.listTrainingRecords)
	mux.HandleFunc("GET /api/training-records/{id}", s.getTrainingRecord)
	mux.HandleFunc("PUT /api/training-records/{id}", s.updateTrainingRecord)
	mux.HandleFunc("DELETE /api/training-records/{id}", s.deleteTrainingRecord)
	mux.HandleFunc("POST /api/training-records/batch-import", s.batchImportTrainingRecords)
	mux.HandleFunc("GET /api/training-records/export", s.exportTrainingRecords)
	mux.HandleFunc("GET /api/training-records/report", s.getTrainingReport)
}

type createTrainingRecordRequest struct {
	CertificateID string `json:"certificate_id"`
	Topic         string `json:"topic"`
	Trainer       string `json:"trainer"`
	Hours         int    `json:"hours"`
	TrainDate     string `json:"train_date"`
	ExamResult    string `json:"exam_result"`
	Status        string `json:"status"`
}

func (s *Server) createTrainingRecord(w http.ResponseWriter, r *http.Request) {
	var req createTrainingRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTrainingRecord(model.TrainingRecord{
		CertificateID: req.CertificateID,
		Topic:         req.Topic,
		Trainer:       req.Trainer,
		Hours:         req.Hours,
		TrainDate:     req.TrainDate,
		ExamResult:    req.ExamResult,
		Status:        req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listTrainingRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TrainingRecordFilter{
		CertificateID: r.URL.Query().Get("certificate_id"),
		Status:        r.URL.Query().Get("status"),
		Topic:         r.URL.Query().Get("topic"),
	}
	items, total, err := s.svc.ListTrainingRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTrainingRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetTrainingRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type updateTrainingRecordRequest struct {
	Topic      string `json:"topic"`
	Trainer    string `json:"trainer"`
	Hours      int    `json:"hours"`
	TrainDate  string `json:"train_date"`
	ExamResult string `json:"exam_result"`
	Status     string `json:"status"`
}

func (s *Server) updateTrainingRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateTrainingRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTrainingRecord(id, model.TrainingRecord{
		Topic:      req.Topic,
		Trainer:    req.Trainer,
		Hours:      req.Hours,
		TrainDate:  req.TrainDate,
		ExamResult: req.ExamResult,
		Status:     req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTrainingRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTrainingRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) batchImportTrainingRecords(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Items []createTrainingRecordRequest `json:"items"`
	}
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.TrainingRecord, len(req.Items))
	for i, item := range req.Items {
		inputs[i] = model.TrainingRecord{
			CertificateID: item.CertificateID,
			Topic:         item.Topic,
			Trainer:       item.Trainer,
			Hours:         item.Hours,
			TrainDate:     item.TrainDate,
			ExamResult:    item.ExamResult,
			Status:        item.Status,
		}
	}
	items, err := s.svc.BatchImportTrainingRecords(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, items)
}

func (s *Server) exportTrainingRecords(w http.ResponseWriter, r *http.Request) {
	data, err := s.svc.ExportTrainingRecordsJSON()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=training_records.json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) getTrainingReport(w http.ResponseWriter, r *http.Request) {
	report, err := s.svc.GetTrainingReport()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, report)
}
