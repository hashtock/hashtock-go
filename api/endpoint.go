package api

import (
    "net/http"
)

type EndPoint struct {
    URI     string
    Method  string
    Name    string
    Handler func(http.ResponseWriter, *http.Request)
}

func (e *EndPoint) isMain() bool {
    return e.URI == "/"
}

func NewEndPoint(uri, method, name string, handler func(http.ResponseWriter, *http.Request)) *EndPoint {
    return &EndPoint{
        URI:     uri,
        Method:  method,
        Name:    name,
        Handler: handler,
    }
}
