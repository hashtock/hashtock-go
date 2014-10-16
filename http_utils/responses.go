package http_utils

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
)

func SerializeResponse(rw http.ResponseWriter, req *http.Request, obj interface{}, status_code int) (err error) {
    accept := req.Header.Get("Accept")

    var data []byte

    q := req.URL.Query()
    switch q.Get("format") {
    case "json":
        accept = "application/json"
    }

    switch accept {
    case "application/json":
        data, err = json.Marshal(obj)
    default:
        err = errors.New(fmt.Sprintf("Unsupported content type: %#v", accept))
    }

    if err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

    rw.Header().Set("Content-Type", accept)
    rw.WriteHeader(status_code)
    rw.Write(data)
    return
}

func SerializeErrorResponse(rw http.ResponseWriter, req *http.Request, err error) error {
    var (
        http_err HttpError
        ok       bool
    )

    if http_err, ok = err.(HttpError); !ok {
        http_err = NewHttpError(http.StatusInternalServerError, err.Error())
    }

    return SerializeResponse(rw, req, http_err, http_err.ErrCode())
}
