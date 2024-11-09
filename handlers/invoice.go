package handlers

import (
	"context"
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
	ctx := context.Background()
	jsonData1, err := h.Cache.Get(ctx, cacheKey).Result()
	if err == nil {
		responses.NewOKResponseWithJson(w, "success", []byte(jsonData1))

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

	jsonData2, err := json.Marshal(response)
	if err != nil {
		return
	}

	err = h.Cache.Set(ctx, cacheKey, jsonData2, time.Minute*5).Err()
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Setting cache failed: %v", err.Error()))
		responses.NewInternalServerError(w, "Server error")
		return
	}

	responses.NewOKResponseWithData(w, "success", data)
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
	ctx := context.Background()
	cacheKey := fmt.Sprintf("invoice_widgets:%v", user.ID)
	_, err = h.Cache.Del(ctx, cacheKey).Result()
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Failed to flush invoice widgets cache: %v", err.Error()))
	}
}
