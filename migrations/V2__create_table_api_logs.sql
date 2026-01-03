CREATE TABLE IF NOT EXISTS api_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    endpoint TEXT NOT NULL,
    http_method VARCHAR(10) NOT NULL,
    status_code INTEGER NOT NULL,
    request_body JSONB,
    response_body JSONB,
    ip_address TEXT,
    user_agent TEXT,
    user_id UUID,
    request_headers JSONB,
    response_headers JSONB,
    duration_ms INTEGER,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);