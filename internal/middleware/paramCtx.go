package middleware

import (
	"context"
	i "github.com/dupreehkuda/balance-microservice/internal"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// ParamCtx extracts accountID from request path
func (m middleware) ParamCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID := chi.URLParam(r, "account_id")
		reportID := chi.URLParam(r, "report_id")

		var (
			akey i.CtxKey = "account"
			rkey i.CtxKey = "report"
		)

		if accountID != "" {
			ctx := context.WithValue(r.Context(), akey, accountID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			ctx := context.WithValue(r.Context(), rkey, reportID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
