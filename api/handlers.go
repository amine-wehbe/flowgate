package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CheckOrigin allows all origins — safe for local dev, lock down before exposing publicly
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Returns all captured requests from DB ordered by most recent first
func getRequests(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := pool.Query(r.Context(), "SELECT * FROM requests ORDER BY captured_at DESC;")
		if err != nil {
			http.Error(w, "DB query error:", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var requests []Request
		for rows.Next() {
			var req Request
			err := rows.Scan(&req.ID, &req.CapturedAt, &req.Method, &req.URL, &req.Host, &req.Path, &req.Protocol, &req.ReqHeaders, &req.ReqBody, &req.ResStatus, &req.ResHeaders, &req.ResBody, &req.DurationMs, &req.TLS)
			if err != nil {
				http.Error(w, "DB request error:", http.StatusInternalServerError)
				return
			}
			requests = append(requests, req)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(requests)
	}
}

// Receives a captured request from the proxy, inserts it into the DB, then pushes it to all WS clients
func postRequests(pool *pgxpool.Pool, hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "DB error:", http.StatusBadRequest)
			return
		}
		_, err = pool.Exec(r.Context(), "INSERT INTO requests (method, url, host, path, protocol, req_headers, req_body, res_status, res_headers, res_body, duration_ms, tls) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", req.Method, req.URL, req.Host, req.Path, req.Protocol, req.ReqHeaders, req.ReqBody, req.ResStatus, req.ResHeaders, req.ResBody, req.DurationMs, req.TLS)
		if err != nil {
			http.Error(w, "DB insertion error:", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		j, err := json.Marshal(req)
		if err != nil {
			log.Println(err)
			return
		}
		hub.broadcast <- j
	}
}

// Deletes all captured requests from the DB
func deleteRequests(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := pool.Exec(r.Context(), "DELETE FROM requests")
		if err != nil {
			http.Error(w, "DB deletion error:", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// Fetches a captured request by ID from the DB and re-fires it, returning the new response
func replayRequests(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		row := pool.QueryRow(r.Context(), "SELECT method, url, req_headers, req_body FROM requests WHERE id = $1", id)
		var req Request
		err := row.Scan(&req.Method, &req.URL, &req.ReqHeaders, &req.ReqBody)
		if err != nil {
			http.Error(w, "DB row scan error:", http.StatusInternalServerError)
			return
		}
		request, err := http.NewRequest(req.Method, req.URL, nil)
		if err != nil {
			http.Error(w, "HTTP request error:", http.StatusInternalServerError)
			return
		}
		for key, values := range req.ReqHeaders {
			request.Header.Set(key, values[0])
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			http.Error(w, "HTTP response error:", http.StatusInternalServerError)
			return
		}
		j, err := io.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			http.Error(w, "Read response error:", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"res_status": response.StatusCode,
			"res_body":   string(j),
		})
	}
}

// Upgrades HTTP to WebSocket, registers with hub, blocks until client disconnects, then unregisters
func serveWs(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Hub serve error:", http.StatusInternalServerError)
			return
		}
		hub.register <- conn
		// Read loop keeps connection alive — ReadMessage errors signal client disconnect
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
		hub.unregister <- conn
	}
}
