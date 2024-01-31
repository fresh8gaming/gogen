package http

import "net/http"

type NotFoundHandler struct {
	Service string
}

func NewNotFoundHandler(service string) NotFoundHandler {
	return NotFoundHandler{
		Service: service,
	}
}

func (nfh NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Service", nfh.Service)
	http.NotFound(w, r)
}
