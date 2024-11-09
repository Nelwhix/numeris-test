package requests

type CreateInvoice struct {
	Title   string `json:"title" validate:"required"`
	Amount  int64  `json:"amount" validate:"required"`
	DueDate string `json:"dueDate" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}
