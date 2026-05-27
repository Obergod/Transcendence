package metrics

import (
    "net/http"
)

func TrackRequests(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        HttpRequests.WithLabelValues(r.URL.Path).Inc()

        next.ServeHTTP(w, r)
    })
}