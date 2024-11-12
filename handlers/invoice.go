package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Nelwhix/numeris/pkg"
	"github.com/Nelwhix/numeris/pkg/models"
	"github.com/Nelwhix/numeris/pkg/requests"
	"github.com/Nelwhix/numeris/pkg/resources"
	"github.com/Nelwhix/numeris/pkg/responses"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) GetInvoiceWidgetsData(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		responses.NewBadRequest(w, "User not found")
		return
	}

	cacheKey := fmt.Sprintf("invoice_widgets:%v", user.ID)
	cachedData, err := h.GetCacheItem(r.Context(), cacheKey)
	if err == nil {
		var response resources.InvoiceWidgetResource
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

	response := resources.InvoiceWidgetResource{
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

	response := resources.InvoiceResource{
		ID:   invoice.ID,
		Type: "invoice",
		Attributes: resources.InvoiceAttributes{
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

	if request.DueDate != "" {
		isValid := pkg.IsValidTime(request.DueDate)
		if !isValid {
			responses.NewBadRequest(w, "Invalid dueAt date")
			return
		}
	}

	var uInvoice models.Invoice
	oInvoice, err := h.Model.GetInvoiceById(r.Context(), invoiceID)
	if err != nil {
		responses.NewNotFound(w, "invoice not found")
		return
	}

	uInvoice.UpdateFromRequest(request, oInvoice)

	h.Logger.Info(fmt.Sprintf("Updating invoice, state is: %v", uInvoice.State.String()))
	err = h.Model.UpdateInvoice(r.Context(), uInvoice)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error updating invoice %v: %v", invoiceID, err.Error()))
		responses.NewInternalServerError(w, err.Error())
	}

	response := resources.InvoiceResource{
		ID:   uInvoice.ID,
		Type: "invoice",
		Attributes: resources.InvoiceAttributes{
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

var AllowedInvoiceSorts = []string{"created_at"}

func (h *Handler) GetInvoices(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		responses.NewBadRequest(w, "User not found")
		return
	}
	limitParam := r.URL.Query().Get("limit")
	limit := 10
	var err error
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			responses.NewBadRequest(w, "limit should be an integer")
			return
		}
	}

	pageParam := r.URL.Query().Get("page")
	page := 1
	if pageParam != "" {
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			responses.NewBadRequest(w, "page should be an integer")
			return
		}
	}

	sortParam := r.URL.Query().Get("sort")
	var allowedSorts []string
	if sortParam != "" {
		userSort := strings.Split(sortParam, ",")

		for _, s := range userSort {
			s = strings.ToLower(strings.TrimSpace(s))
			present := slices.Contains(AllowedInvoiceSorts, s[1:])
			if present {
				allowedSorts = append(allowedSorts, s)
			}
		}
	}

	invoices, err := h.Model.GetInvoices(r.Context(), user.ID, limit, page, allowedSorts)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error getting invoices for user with id: %v %v", user.ID, err.Error()))
		responses.NewInternalServerError(w, "Server error")
		return
	}

	invoiceCount, err := h.Model.GetTotalInvoiceCountForUser(r.Context(), user.ID)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Error getting total number of invoices for user with id: %v %v", user.ID, err.Error()))
		responses.NewInternalServerError(w, "Server error")
		return
	}

	var totalPages int
	if invoiceCount.Total < limit {
		totalPages = 1
	} else {
		totalPages = invoiceCount.Total / limit
	}

	response := resources.InvoiceListResource{
		Message: "Get Invoices.",
		Data:    resources.FromInvoices(invoices),
		Meta: resources.Meta{
			Pagination: resources.Pagination{
				CurrentPage: page,
				PerPage:     limit,
				Count:       len(invoices),
				Total:       invoiceCount.Total,
				TotalPages:  totalPages,
			},
		},
	}

	responses.NewOK(w, response)
}

func (h *Handler) GetInvoiceActivities(w http.ResponseWriter, r *http.Request) {
	responses.NewOK(w, "implement me")
	//user, ok := r.Context().Value("user").(models.User)
	//if !ok {
	//	responses.NewBadRequest(w, "User not found")
	//	return
	//}
	//
	//invoices, err := h.Model.GetInvoices(r.Context(), user.ID, limit, page, allowedSorts)
	//if err != nil {
	//	h.Logger.Error(fmt.Sprintf("Error getting invoices for user with id: %v %v", user.ID, err.Error()))
	//	responses.NewInternalServerError(w, "Server error")
	//	return
	//}
	//
	//response := resources.InvoiceListResource{
	//	Message: "Get Invoices.",
	//	Data:    resources.FromInvoices(invoices),
	//	Meta: resources.Meta{
	//		Pagination: resources.Pagination{
	//			CurrentPage: page,
	//			PerPage:     limit,
	//			Count:       len(invoices),
	//			Total:       invoiceCount.Total,
	//			TotalPages:  totalPages,
	//		},
	//	},
	//}
	//
	//responses.NewOK(w, response)
}
