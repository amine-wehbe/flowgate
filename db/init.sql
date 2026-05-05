CREATE TABLE IF NOT EXISTS requests (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  captured_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  method      TEXT NOT NULL,
  url         TEXT NOT NULL,
  host        TEXT NOT NULL,
  path        TEXT NOT NULL,
  protocol    TEXT NOT NULL,
  req_headers JSONB,
  req_body    TEXT,
  res_status  INTEGER,
  res_headers JSONB,
  res_body    TEXT,
  duration_ms INTEGER,
  tls         BOOLEAN DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_requests_host ON requests(host);
CREATE INDEX IF NOT EXISTS idx_requests_captured_at ON requests(captured_at DESC);
