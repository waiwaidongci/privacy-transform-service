CREATE TABLE IF NOT EXISTS privacy_classifications (id TEXT PRIMARY KEY, name TEXT NOT NULL, level INT NOT NULL, description TEXT);
CREATE TABLE IF NOT EXISTS privacy_policies (id TEXT NOT NULL, version INT NOT NULL, name TEXT NOT NULL, status TEXT NOT NULL, rules JSONB NOT NULL, updated_at TIMESTAMPTZ NOT NULL, PRIMARY KEY(id,version));
CREATE TABLE IF NOT EXISTS privacy_token_metadata (token TEXT PRIMARY KEY, policy_id TEXT NOT NULL, key_version TEXT NOT NULL, revoked BOOLEAN NOT NULL DEFAULT FALSE, created_at TIMESTAMPTZ NOT NULL);
CREATE TABLE IF NOT EXISTS privacy_transform_requests (id TEXT PRIMARY KEY, policy_id TEXT NOT NULL, summary JSONB NOT NULL, created_at TIMESTAMPTZ NOT NULL);
CREATE TABLE IF NOT EXISTS privacy_transform_results (request_id TEXT PRIMARY KEY, policy_id TEXT NOT NULL, result_hash TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL);
CREATE INDEX IF NOT EXISTS idx_privacy_requests_policy ON privacy_transform_requests(policy_id,created_at);
