package pkg

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"\n{\n\tMethod: %s,\n\tRequestURI: %s,\n\tTime: %s,\n}\n",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
