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

    switch accept {
    case "application/json":
        data, err = json.Marshal(obj)
    default:
        err = errors.New(fmt.Sprintf("Unsupported content type: %s", accept))
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
