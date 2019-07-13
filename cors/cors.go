package cors

import (
	"net/http"
)

var PermittedOrigins = map[string]struct{}{
	"http://localhost:3000": struct{}{},
}

// Middleware correctly performs Cross Origin Checks
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("origin")

		if _, ok := PermittedOrigins[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, X-Requested-With, User-IP, Referer, Lomas") //  Cache-Control,
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if req.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400") // 1 day
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Cache-Control", "no-cache")
		h.ServeHTTP(w, req)
	})
}
