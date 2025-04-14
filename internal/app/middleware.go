package app

import (
	"avito/internal/utils"
	"avito/internal/utils/errors"
	"avito/internal/utils/jwt"
	"avito/internal/utils/roles"
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
			errors.Write(w, errors.ErrInvalidToken, http.StatusForbidden)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.Parse(tokenStr)
		if err != nil {
			errors.Write(w, errors.ErrInvalidToken, http.StatusForbidden)
			return
		}

		ctx := roles.WithRole(r.Context(), claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
