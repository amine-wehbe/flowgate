package main

import (
	"time"
)

// Request mirrors the requests table — used for DB scanning and JSON serialization
type Request struct {
	ID         string              `json:"id"`
	CapturedAt time.Time           `json:"captured_at"`
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Host       string              `json:"host"`
	Path       string              `json:"path"`
	Protocol   string              `json:"protocol"`
	ReqHeaders map[string][]string `json:"req_headers"`
	ReqBody    string              `json:"req_body"`
	ResStatus  int                 `json:"res_status"`
	ResHeaders map[string][]string `json:"res_headers"`
	ResBody    string              `json:"res_body"`
	DurationMs int                 `json:"duration_ms"`
	TLS        bool                `json:"tls"`
}
