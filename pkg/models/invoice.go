package models

import (
	"context"
	"github.com/Nelwhix/numeris/pkg/enums"
	"github.com/oklog/ulid/v2"
	"time"
)

type Invoice struct {
	ID        string
	Title     string
	UserID    string
	Amount    int64
	State     enums.InvoiceState
	DueAt     time.Time
	CreatedAt time.Time
}

type InvoiceWidget struct {
	TotalDraft         *int64
	TotalOverdue       *int64
	TotalUnpaid        *int64
	TotalPaid          *int64
	TotalDraftAmount   *int64
	TotalOverdueAmount *int64
	TotalUnpaidAmount  *int64
	TotalPaidAmount    *int64
}

func (m *Model) InsertIntoInvoices(ctx context.Context, invoice Invoice) (Invoice, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	invoiceID := ulid.Make().String()

	sql := "insert into invoices(id, title, amount_cents, state, due_at, user_id) values ($1, $2, $3, $4, $5, $6)"
	_, err := m.Conn.Exec(ctx, sql, invoiceID, invoice.Title, invoice.Amount, enums.Draft, invoice.DueAt, invoice.UserID)

	if err != nil {
		return Invoice{}, err
	}

	invoice, err = m.GetInvoiceById(ctx, invoiceID)
	if err != nil {
		return Invoice{}, err
	}

	return invoice, nil
}

func (m *Model) GetInvoiceById(ctx context.Context, invoiceID string) (Invoice, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var invoice Invoice
	var state int

	row := m.Conn.QueryRow(ctx, "select id, title, amount_cents, state, due_at, created_at FROM invoices WHERE id = $1", invoiceID)
	err := row.Scan(&invoice.ID, &invoice.Title, &invoice.Amount, &state, &invoice.DueAt, &invoice.CreatedAt)
	if err != nil {
		return Invoice{}, err
	}

	invoice.State, err = enums.ParseInvoiceState(state)
	if err != nil {
		return Invoice{}, err
	}

	return invoice, nil
}

func (m *Model) GetInvoiceWidgetsData(ctx context.Context, userID string) (InvoiceWidget, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	sql := `SELECT
    	COUNT(CASE WHEN state = 0 THEN 1 END) AS total_draft_invoices,
    	COUNT(CASE WHEN state = 1 THEN 1 END) AS total_overdue_invoices,
		COUNT(CASE WHEN state = 2 THEN 1 END) AS total_unpaid_invoices,
		COUNT(CASE WHEN state = 3 THEN 1 END) AS total_paid_invoices,
		SUM(CASE WHEN state = 0 THEN amount_cents END) AS total_draft_amount,
		SUM(CASE WHEN state = 1 THEN amount_cents END) AS total_overdue_amount,
		SUM(CASE WHEN state = 2 THEN amount_cents END) AS total_unpaid_amount,
		SUM(CASE WHEN state = 3 THEN amount_cents END) AS total_paid_amount
	FROM invoices where user_id = $1`

	var invoiceWidget InvoiceWidget
	row := m.Conn.QueryRow(ctx, sql, userID)
	err := row.Scan(&invoiceWidget.TotalDraft, &invoiceWidget.TotalOverdue, &invoiceWidget.TotalUnpaid,
		&invoiceWidget.TotalPaid, &invoiceWidget.TotalDraftAmount, &invoiceWidget.TotalOverdueAmount,
		&invoiceWidget.TotalUnpaidAmount, &invoiceWidget.TotalPaidAmount)

	if err != nil {
		return InvoiceWidget{}, err
	}

	return invoiceWidget, nil
}
