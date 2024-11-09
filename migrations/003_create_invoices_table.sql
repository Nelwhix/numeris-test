CREATE TABLE IF NOT EXISTS invoices (
     id CHAR(26) PRIMARY KEY,
     user_id CHAR(26) NOT NULL,
     title VARCHAR(255) NOT NULL,
     amount_cents bigint NOT NULL,
     state smallint not null,
     due_at TIMESTAMP NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invoices_user_id ON invoices (user_id);
CREATE INDEX idx_invoices_state ON invoices (state);
CREATE INDEX idx_invoices_due_at ON invoices (due_at);
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
