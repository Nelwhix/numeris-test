CREATE TABLE IF NOT EXISTS invoices (
     id CHAR(26) PRIMARY KEY,
     title VARCHAR(255) NOT NULL,
     amount bigint NOT NULL,
     due_at TIMESTAMP NOT NULL,
     state smallint not null,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invoices_state ON invoices (state);
CREATE INDEX idx_invoices_due_at ON invoices (due_at);
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
