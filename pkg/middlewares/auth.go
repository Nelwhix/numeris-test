package middlewares

import (
	"context"
	"github.com/Nelwhix/numeris/pkg/models"
	"github.com/Nelwhix/numeris/pkg/responses"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	Model models.Model
}

func (a *AuthMiddleware) Register(handlerFunc func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			responses.NewUnauthorized(w, "Unauthorized.")
			return
		}

		user, err := a.Model.GetUserByToken(r.Context(), parts[1])
		if err != nil {
			responses.NewUnauthorized(w, "Unauthorized.")
			return
		}

		nContext := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(nContext)

		http.HandlerFunc(handlerFunc).ServeHTTP(w, r)
	})
}
