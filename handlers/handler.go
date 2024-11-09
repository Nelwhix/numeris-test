package handlers

import (
	"github.com/Nelwhix/numeris/pkg/models"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type Handler struct {
	Model     models.Model
	Logger    *slog.Logger
	Validator *validator.Validate
	Cache     *redis.Client
}
