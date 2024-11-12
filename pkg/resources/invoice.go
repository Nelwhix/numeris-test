package resources

import (
	"github.com/Nelwhix/numeris/pkg"
	"github.com/Nelwhix/numeris/pkg/enums"
	"github.com/Nelwhix/numeris/pkg/models"
	"time"
)

type InvoiceResource struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes InvoiceAttributes `json:"attributes"`
}

type InvoiceAttributes struct {
	Title     string             `json:"title"`
	Amount    string             `json:"amount"`
	State     enums.InvoiceState `json:"state"`
	DueAt     time.Time          `json:"dueAt"`
	CreatedAt time.Time          `json:"createdAt"`
}

type InvoiceWidgetResource struct {
	TotalDraft         *int64 `json:"totalDraft"`
	TotalOverdue       *int64 `json:"totalOverdue"`
	TotalUnpaid        *int64 `json:"totalUnpaid"`
	TotalPaid          *int64 `json:"totalPaid"`
	TotalDraftAmount   string `json:"totalDraftAmount"`
	TotalOverdueAmount string `json:"totalOverdueAmount"`
	TotalUnpaidAmount  string `json:"totalUnpaidAmount"`
	TotalPaidAmount    string `json:"totalPaidAmount"`
}

type InvoiceListResource struct {
	Message string            `json:"message"`
	Data    []InvoiceResource `json:"data"`
	Meta    Meta              `json:"meta"`
}

type Meta struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	PerPage     int `json:"perPage"`
	Count       int `json:"count"`
	Total       int `json:"total"`
	TotalPages  int `json:"totalPages"`
}

func FromInvoices(invoices []models.Invoice) []InvoiceResource {
	var invoiceRes []InvoiceResource

	for _, invoice := range invoices {
		i := InvoiceResource{
			ID:   invoice.ID,
			Type: "invoice",
			Attributes: InvoiceAttributes{
				Title:     invoice.Title,
				Amount:    pkg.FormatMoneyToUsd(&invoice.Amount),
				State:     invoice.State,
				DueAt:     invoice.DueAt,
				CreatedAt: invoice.CreatedAt,
			},
		}

		invoiceRes = append(invoiceRes, i)
	}

	return invoiceRes
}
