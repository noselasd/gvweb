package httputil

import (
	"fmt"
	"log"
	"net/http"
)

func accessLogger(l chan string) {
	for {
		log.Print(<-l)
	}
}

type statusWrapper struct {
	status int
	http.ResponseWriter
}

func (s *statusWrapper) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

func httpLog(l chan string, sw *statusWrapper, r *http.Request) {

	var remote string
	if len(r.Header["X-Forwarded-For"]) > 0 {
		remote = r.Header["X-Forwarded-For"][0]
	} else {
		remote = r.RemoteAddr
	}
	l <- fmt.Sprintf("%s %s %s %d", remote, r.Method, r.URL, sw.status)
}

func NewLogWrapper(handler http.Handler) http.Handler {
	lchan := make(chan string, 64)
	go accessLogger(lchan)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := statusWrapper{-1, w}
		w.Header().Set("Server", "Go HTTP handler")
		handler.ServeHTTP(&sw, r)
		httpLog(lchan, &sw, r)
	})
}
