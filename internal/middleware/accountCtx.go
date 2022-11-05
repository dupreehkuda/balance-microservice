package middleware

import (
	"context"
	i "github.com/dupreehkuda/balance-microservice/internal"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// AccountCtx extracts accountID from request path
func (m middleware) AccountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID := chi.URLParam(r, "articleID")

		var key i.CtxKey = "account"

		ctx := context.WithValue(r.Context(), key, accountID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
