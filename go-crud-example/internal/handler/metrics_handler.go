package handler

import (
    "github.com/gorilla/mux"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetricsHandlers(r *mux.Router) {
    r.Handle("/metrics", promhttp.Handler())
}
