package metrics

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GetRouter() http.Handler {
	r := mux.NewRouter()

	r.Handle("/_metrics", promhttp.Handler())

	return handlers.CompressHandler(r)
}
