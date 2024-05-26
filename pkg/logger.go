package pkg

import (
	"log"
	"net/http"
	"time"
)

type customResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.status = code
	crw.ResponseWriter.WriteHeader(code)
}

func (crw *customResponseWriter) Write(data []byte) (int, error) {
	if crw.status == 0 {
		crw.status = http.StatusOK
	}
	size, err := crw.ResponseWriter.Write(data)
	crw.size += size
	return size, err
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		crw := &customResponseWriter{w, http.StatusOK, 0}

		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Panic: %v", err)
			}

			remoteIP := r.RemoteAddr
			host := r.Host
			uri := r.RequestURI
			method := r.Method

			duration := time.Since(start)

			log.Printf(
				"\n{\n\tMethod: %s,\n\tRequestURI: %s,\n\tStatus: %d,\n\tTime: %s,\n\tremote_ip: %s,\n\thost: %s,\n\turi: %s,\n\tmethod: %s,\n\tuser_agent: %s,\n\tbytes_out: %d,\n}\n",
				method,
				uri,
				crw.status,
				duration,
				remoteIP,
				host,
				uri,
				method,
				r.UserAgent(),
				crw.size,
			)
		}()

		next.ServeHTTP(crw, r)
	})
}
