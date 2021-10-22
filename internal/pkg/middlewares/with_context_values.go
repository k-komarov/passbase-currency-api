package middlewares

import (
	"context"
	"net/http"
)

func WithContextValues(keyValues map[interface{}]interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			for key, value := range keyValues {
				ctx = context.WithValue(ctx, key, value)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
