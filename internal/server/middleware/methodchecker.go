package middleware

import "net/http"

func MethodChecker(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		}
	})
}
