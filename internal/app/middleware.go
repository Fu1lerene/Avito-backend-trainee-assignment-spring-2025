package app

import (
	"avito/internal/api"
	"avito/internal/utils"
	"avito/internal/utils/jwt"
	"avito/internal/utils/roles"
	"log/slog"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger/") || utils.IsPathSkipped(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			api.WriteForbidden(w, "invalid token or missing token")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.Parse(tokenStr)
		if err != nil {
			slog.Error("token parse error", "error", err.Error())
			api.WriteForbidden(w, "invalid token")
			return
		}

		ctx := roles.WithRole(r.Context(), claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
