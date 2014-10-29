package http_utils

import (
    "net/http"
    "time"

    "appengine"

    "github.com/codegangsta/negroni"
)

type RequestTimer struct{}

func NewRequestTimer() *RequestTimer {
    return &RequestTimer{}
}

func (rt *RequestTimer) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    cxt := appengine.NewContext(r)
    start := time.Now()

    next(rw, r)
    res := rw.(negroni.ResponseWriter)
    cxt.Debugf("Request \"%v %v\" completed with %v in %.2fms", r.Method, r.URL, res.Status(), rt.time_since_ms(start))
}

func (rt *RequestTimer) time_since_ms(start time.Time) float64 {
    return time.Since(start).Seconds() * 1000
}
