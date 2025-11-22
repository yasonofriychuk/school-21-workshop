package auth

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

const userDataKey = contextKey("userData")

type Middleware struct{}

func NewAuthMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	// Заглушка для workshop
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(r.Header.Get("Authorization"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userDataKey, userId)))
	})
}

// GetUserId - вспомогательная функция для извлечения данных из контекста.
func GetUserId(ctx context.Context) (int64, bool) {
	data, ok := ctx.Value(userDataKey).(int64)
	return data, ok
}
