# Flowgate

A self-hosted HTTP/HTTPS intercepting proxy with a live inspection UI. Think simplified Burp Suite + Charles Proxy, built from scratch.

## Stack

- **Proxy** — Go, handles HTTP and HTTPS (TLS MITM via local CA)
- **API** — Go, REST + WebSocket
- **Database** — PostgreSQL
- **Frontend** — React + Vite, served by nginx in Docker

## Features

- Intercept HTTP and HTTPS traffic transparently
- Live request feed via WebSocket — new requests appear instantly
- Full request/response detail: headers, body, status, timing, TLS flag
- Replay any captured request directly from the UI
- Clear history with one click
- Dark theme, two-panel layout

## Quick Start

**Prerequisites:** Docker Desktop, mkcert installed and CA trusted (`mkcert -install`)

```bash
git clone https://github.com/amine-wehbe/flowgate
cd flowgate
docker-compose up --build
```

- Frontend: http://localhost:4000
- API: http://localhost:3000
- Proxy: localhost:8080

## Configure Your Browser

To capture browser traffic, set your system proxy to `localhost:8080` for both HTTP and HTTPS:

**Mac:** System Settings → Network → Wi-Fi → Details → Proxies
- Web Proxy (HTTP): `localhost:8080`
- Secure Web Proxy (HTTPS): `localhost:8080`

Remember to disable the proxy when done.

## Configure mkcert Volume

The proxy container needs access to your mkcert CA to do TLS MITM. Update the volume path in `docker-compose.yml` to match your machine:

```yaml
volumes:
  - /path/to/your/mkcert:/certs:ro
```

Find your mkcert path with: `mkcert -CAROOT`

## Reset Database

```bash
docker-compose down -v && docker-compose up
```

## Local Development (without Docker)

```bash
# Terminal 1 — DB only
docker-compose up db

# Terminal 2 — API
cd api && go run main.go db.go models.go handlers.go hub.go

# Terminal 3 — Proxy
cd proxy && go run main.go proxy.go tls.go sender.go

# Terminal 4 — Frontend
cd frontend && npm run dev
```

Frontend at http://localhost:5173

## Architecture

```
Browser → Proxy (:8080) → API (:3000) → PostgreSQL
                                ↓
                         WebSocket hub
                                ↓
                       React UI (:4000)
```

The proxy intercepts all traffic, logs it to the API, which persists to PostgreSQL and broadcasts to all connected WebSocket clients in real time.
