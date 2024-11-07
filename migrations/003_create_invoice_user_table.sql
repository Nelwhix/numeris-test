CREATE TABLE IF NOT EXISTS invoice_user (
      user_id CHAR(26) NOT NULL,
      total_paid_cents bigint DEFAULT 0,
      total_overdue_cents bigint DEFAULT 0,
      total_draft_cents bigint DEFAULT 0,
      total_unpaid_cents bigint DEFAULT 0,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invoice_user_user_id ON invoice_user (user_id);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
