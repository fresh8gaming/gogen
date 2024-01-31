package middleware

import (
	"net/http"
)

func Service(service string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Service", service)
			next.ServeHTTP(w, r)
		})
	}
}
