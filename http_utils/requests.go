package http_utils

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
)

func DeSerializeRequest(req http.Request, obj interface{}) (err error) {
    content_type := req.Header.Get("content-type")

    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        return err
    }

    switch content_type {
    case "application/json":
        err = json.Unmarshal(body, &obj)
    default:
        err = errors.New(fmt.Sprintf("Unsupported content type: %s", content_type))
    }
    return
}
