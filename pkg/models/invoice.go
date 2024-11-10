package models

import (
	"context"
	"github.com/Nelwhix/numeris/pkg/enums"
	"github.com/Nelwhix/numeris/pkg/requests"
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

func (i *Invoice) UpdateFromRequest(request requests.UpdateInvoice, o Invoice) {
	if request.Title == "" {
		i.Title = o.Title
	} else {
		i.Title = request.Title
	}

	if request.Amount == 0 {
		i.Amount = o.Amount
	} else {
		i.Amount = request.Amount
	}

	if request.State == nil {
		i.State = o.State
	} else {
		i.State = *request.State
	}

	if request.DueDate == "" {
		i.DueAt = o.DueAt
	} else {
		i.DueAt, _ = time.Parse(time.RFC3339, request.DueDate)
	}
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

func (m *Model) UpdateInvoice(ctx context.Context, invoice Invoice) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	sql := "update invoices set title = $1, amount_cents = $2, state = $3, due_at = $4 where id = $5"
	_, err := m.Conn.Exec(ctx, sql, invoice.Title, invoice.Amount, invoice.State, invoice.DueAt, invoice.ID)
	if err != nil {
		return err
	}

	return nil
}
