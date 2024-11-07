package handlers

import (
	"github.com/Nelwhix/numeris/pkg/models"
	"github.com/go-playground/validator/v10"
	"log/slog"
)

type Handler struct {
	Model     models.Model
	Logger    *slog.Logger
	Validator *validator.Validate
}
