CREATE TABLE IF NOT EXISTS personal_access_tokens (
      user_id CHAR(26) NOT NULL,
      token VARCHAR(64) NOT NULL,
      last_used_at TIMESTAMP null,
      expires_at TIMESTAMP NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_token ON personal_access_tokens (token);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
