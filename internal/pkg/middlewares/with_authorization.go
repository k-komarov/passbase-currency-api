package middlewares

import (
	"context"
	"github.com/k-komarov/passbase-currency-api/internal/pkg/constants"
	"github.com/k-komarov/passbase-currency-api/internal/pkg/model"
	"net/http"
	"strings"
)

func WithAuthorization(projects map[string]*model.Project) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
			ctx := r.Context()
			if authToken == "" {
				next.ServeHTTP(w, r)
				return
			}
			if p, ok := projects[authToken]; ok {
				ctx = context.WithValue(ctx, constants.CTX_PROJECT, p)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
