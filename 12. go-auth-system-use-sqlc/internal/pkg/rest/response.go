package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"go-auth-system/internal/util"
)

// APIResponse represents a unified API response structure.
// @Description Unified API response with generics support
type APIResponse[T any] struct {
	IsSuccess bool        `json:"is_success"`
	Message   string      `json:"message,omitempty"`
	Data      T           `json:"data,omitempty"`
	Errors    []ErrorItem `json:"errors,omitempty"` // Error details, business rule
	Meta      *Meta       `json:"meta,omitempty"`   // Pagination, Request ID
}

type ErrorItem struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

type Meta struct {
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
	Duration  string    `json:"duration"`
}

// Success trả về response thành công (200 OK)
func Success[T any](w http.ResponseWriter, r *http.Request, data T, message string) {
	send(w, r, http.StatusOK, true, data, message, nil)
}

// Created trả về response tạo mới thành công (201 Created)
func Created[T any](w http.ResponseWriter, r *http.Request, data T, message string) {
	send(w, r, http.StatusCreated, true, data, message, nil)
}

// Error trả về lỗi chung (4xx, 5xx)
func Error(w http.ResponseWriter, r *http.Request, status int, message string, errs []ErrorItem) {
	send[any](w, r, status, false, nil, message, errs)
}

// Internal helper để xử lý logic chung
func send[T any](w http.ResponseWriter, r *http.Request, status int, isSuccess bool, data T, message string, errs []ErrorItem) {
	reqID, _ := r.Context().Value(util.RequestIDKey).(string)
	startTime, _ := r.Context().Value(util.StartTimeKey).(time.Time)

	res := APIResponse[T]{
		IsSuccess: isSuccess,
		Message:   message,
		Data:      data,
		Errors:    errs,
		Meta: &Meta{
			RequestID: reqID,
			Timestamp: time.Now(),
			Duration:  time.Since(startTime).String(), // Tính toán latency
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}
