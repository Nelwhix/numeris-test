package requests

import "github.com/Nelwhix/numeris/pkg/enums"

type CreateInvoice struct {
	Title   string `json:"title" validate:"required"`
	Amount  int64  `json:"amount" validate:"required"`
	DueDate string `json:"dueDate" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

type UpdateInvoice struct {
	Title   string              `json:"title"`
	Amount  int64               `json:"amount"`
	State   *enums.InvoiceState `json:"state"`
	DueDate string              `json:"dueDate" validate:"datetime=2006-01-02T15:04:05Z07:00"`
}
