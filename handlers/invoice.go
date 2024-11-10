package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Nelwhix/numeris/pkg"
	"github.com/Nelwhix/numeris/pkg/enums"
	"github.com/Nelwhix/numeris/pkg/models"
	"github.com/Nelwhix/numeris/pkg/requests"
	"github.com/Nelwhix/numeris/pkg/responses"
	"net/http"
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

func (h *Handler) GetInvoiceWidgetsData(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		responses.NewBadRequest(w, "User not found")
		return
	}

	cacheKey := fmt.Sprintf("invoice_widgets:%v", user.ID)
	cachedData, err := h.GetCacheItem(r.Context(), cacheKey)
	if err == nil {
		var response InvoiceWidgetResource
		err := json.Unmarshal(cachedData, &response)
		if err != nil {
			return
		}

		responses.NewOKResponseWithData(w, "success", response)
		return
	}

	data, err := h.Model.GetInvoiceWidgetsData(r.Context(), user.ID)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Failed to get invoice widgets data: %v", err.Error()))
		responses.NewInternalServerError(w, "Server error")
		return
	}

	response := InvoiceWidgetResource{
		TotalDraft:         data.TotalDraft,
		TotalOverdue:       data.TotalOverdue,
		TotalUnpaid:        data.TotalUnpaid,
		TotalPaid:          data.TotalPaid,
		TotalDraftAmount:   pkg.FormatMoneyToUsd(data.TotalDraftAmount),
		TotalOverdueAmount: pkg.FormatMoneyToUsd(data.TotalOverdueAmount),
		TotalUnpaidAmount:  pkg.FormatMoneyToUsd(data.TotalUnpaidAmount),
		TotalPaidAmount:    pkg.FormatMoneyToUsd(data.TotalPaidAmount),
	}

	responses.NewOKResponseWithData(w, "success", response)

	err = h.SetCacheItem(r.Context(), cacheKey, response)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Setting invoice widgets cache failed: %v", err.Error()))
		return
	}
}

func (h *Handler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	request, err := pkg.ParseRequestBody[requests.CreateInvoice](r)
	if err != nil {
		responses.NewUnprocessableEntity(w, err.Error())
		return
	}
	err = h.Validator.Struct(request)
	if err != nil {
		responses.NewUnprocessableEntity(w, err.Error())
		return
	}

	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		responses.NewBadRequest(w, "User not found")
		return
	}

	dueAtTime, err := time.Parse(time.RFC3339, request.DueDate)
	if err != nil {
		responses.NewBadRequest(w, "Invalid dueAt date")
		return
	}

	invoice := models.Invoice{
		Title:  request.Title,
		Amount: request.Amount,
		DueAt:  dueAtTime,
		UserID: user.ID,
	}
	invoice, err = h.Model.InsertIntoInvoices(r.Context(), invoice)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error creating invoice: %v", err.Error()))
		responses.NewInternalServerError(w, "Server error")
		return
	}

	response := InvoiceResource{
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

	responses.NewCreatedResponseWithData(w, "Invoice created successfully.", response)

	cacheKey := fmt.Sprintf("invoice_widgets:%v", user.ID)
	err = h.DeleteCacheItem(r.Context(), cacheKey)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Failed to flush invoice widgets cache: %v", err.Error()))
	}
}

func (h *Handler) UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		responses.NewBadRequest(w, "User not found")
		return
	}

	invoiceID := r.PathValue("invoiceID")
	request, err := pkg.ParseRequestBody[requests.UpdateInvoice](r)
	if err != nil {
		responses.NewUnprocessableEntity(w, err.Error())
		return
	}

	isValid := pkg.IsValidTime(request.DueDate)
	if !isValid {
		responses.NewBadRequest(w, "Invalid dueAt date")
		return
	}

	var uInvoice models.Invoice
	oInvoice, err := h.Model.GetInvoiceById(r.Context(), invoiceID)
	if err != nil {
		responses.NewNotFound(w, "invoice not found")
		return
	}

	uInvoice.UpdateFromRequest(request, oInvoice)

	err = h.Model.UpdateInvoice(r.Context(), uInvoice)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error updating invoice %v: %v", invoiceID, err.Error()))
		responses.NewInternalServerError(w, err.Error())
	}

	response := InvoiceResource{
		ID:   uInvoice.ID,
		Type: "invoice",
		Attributes: InvoiceAttributes{
			Title:     uInvoice.Title,
			Amount:    pkg.FormatMoneyToUsd(&uInvoice.Amount),
			State:     uInvoice.State,
			DueAt:     uInvoice.DueAt,
			CreatedAt: uInvoice.CreatedAt,
		},
	}

	responses.NewOKResponseWithData(w, "Invoice updated successfully.", response)

	cacheKey := fmt.Sprintf("invoice_widgets:%v", user.ID)
	err = h.DeleteCacheItem(r.Context(), cacheKey)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Failed to flush invoice widgets cache: %v", err.Error()))
	}
}
