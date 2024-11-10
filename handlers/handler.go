package handlers

import (
	"context"
	"encoding/json"
	"github.com/Nelwhix/numeris/pkg/models"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

type Handler struct {
	Model     models.Model
	Logger    *slog.Logger
	Validator *validator.Validate
	Cache     *redis.Client
}

func (h *Handler) GetCacheItem(ctx context.Context, cacheKey string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	result, err := h.Cache.Get(ctx, cacheKey).Result()
	if err != nil {
		return []byte(""), err
	}

	return []byte(result), nil
}

func (h *Handler) SetCacheItem(ctx context.Context, cacheKey string, data interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	jsonResponse, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = h.Cache.Set(ctx, cacheKey, jsonResponse, time.Minute*5).Err()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) DeleteCacheItem(ctx context.Context, cacheKey string) error {
	_, err := h.Cache.Del(ctx, cacheKey).Result()
	if err != nil {
		return err
	}

	return nil
}
