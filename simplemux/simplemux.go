package simplemux

import (
	"net/http"
	"regexp"
)

type route struct {
	re      *regexp.Regexp
	handler func(http.ResponseWriter, *http.Request, []string)
}

type RegexpHandler struct {
	routes          []*route
	NotFoundHandler http.Handler
}

func NewRegexpHandler() *RegexpHandler {
	h := RegexpHandler{}
	h.NotFoundHandler = http.NotFoundHandler()

	return &h
}

func (h *RegexpHandler) AddRoute(re string, handler func(http.ResponseWriter, *http.Request, []string)) {
	r := &route{regexp.MustCompile(re), handler}
	h.routes = append(h.routes, r)
}

func (h *RegexpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		matches := route.re.FindStringSubmatch(r.URL.Path)
		if matches != nil {
			route.handler(rw, r, matches)
			return
		}
	}

	h.NotFoundHandler.ServeHTTP(rw, r)
}
