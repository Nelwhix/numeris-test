CREATE TABLE IF NOT EXISTS invoice_activities (
    id CHAR(26) PRIMARY KEY,
    user_id CHAR(26) NOT NULL,
    invoice_id CHAR(26) NOT NULL,
    event VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invoice_activities_invoice_id ON invoice_activities (invoice_id);
CREATE INDEX idx_invoice_activities_user_id ON invoice_activities (user_id);