package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CapturedRequest struct {
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

// Sends a captured request/response pair to the API for persistence
func sendToAPI(req *http.Request, res *http.Response, duration time.Duration, isTLS bool) {
	reqBody, _ := io.ReadAll(req.Body)
	resBody, _ := io.ReadAll(res.Body)
	// Build the payload from the intercepted request and response
	r := CapturedRequest{
		Method:     req.Method,
		URL:        req.URL.String(),
		Host:       req.Host,
		Path:       req.URL.Path,
		Protocol:   req.Proto,
		ReqHeaders: req.Header,
		ReqBody:    string(reqBody),
		ResStatus:  res.StatusCode,
		ResHeaders: res.Header,
		ResBody:    string(resBody),
		DurationMs: int(duration.Milliseconds()),
		TLS:        isTLS,
	}
	body, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
		return
	}

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3000"
	}
	_, err = http.Post(apiURL+"/api/requests", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		return
	}
}
