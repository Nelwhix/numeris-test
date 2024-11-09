package middlewares

import (
	"context"
	"github.com/Nelwhix/numeris/pkg/models"
	"github.com/Nelwhix/numeris/pkg/responses"
	"log/slog"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	Model  models.Model
	Logger *slog.Logger
}

func (a *AuthMiddleware) Register(handlerFunc func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			a.Logger.Info("Parts couldn't be split", "Authorization", authHeader)
			responses.NewUnauthorized(w, "Unauthorized.")
			return
		}

		user, err := a.Model.GetUserByToken(r.Context(), parts[1])
		if err != nil {
			a.Logger.Info("User couldn't be found")
			responses.NewUnauthorized(w, "Unauthorized.")
			return
		}

		nContext := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(nContext)

		http.HandlerFunc(handlerFunc).ServeHTTP(w, r)
	})
}

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
