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

		remoteIP := r.RemoteAddr
		host := r.Host
		uri := r.RequestURI
		method := r.Method

		log.Printf(
			"\n{\n\tMethod: %s,\n\tRequestURI: %s,\n\tTime: %s,\n\tremote_ip: %s,\n\thost: %s,\n\turi: %s,\n\tmethod: %s,\n\tuser_agent: %s,\n}\n",
			method,
			uri,
			time.Since(start),
			remoteIP,
			host,
			uri,
			method,
			r.UserAgent(),
		)
	})
}
