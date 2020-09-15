package urlshorten

import (
	"log"
	"net/http"
	"time"
)

type Middleware struct {
}

func (m *Middleware) LoggingHandler(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()

		log.Println("[LoggingMiddleware] URI:", r.RequestURI)
		log.Println("[LoggingMiddleware] TimeCost:", t2.Sub(t1))
	}
	return http.HandlerFunc(f)
}
