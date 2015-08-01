package serialize

import (
	"encoding/json"
	"log"
	"net/http"
)

type Serializer interface {
	JSON(rw http.ResponseWriter, status int, obj interface{})
}

type WebAPISerializer struct{}

func (w WebAPISerializer) JSON(rw http.ResponseWriter, status int, obj interface{}) {
	if err, ok := obj.(error); ok {
		obj = err.Error()
	}

	data, err := json.Marshal(obj)
	if err != nil {
		msg := "Could not serialize object to JSON"
		http.Error(rw, msg, http.StatusInternalServerError)
		log.Printf("%v. Obj: %#v. Err: %v", msg, obj, err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)

	if obj != nil && status != http.StatusNoContent {
		rw.Write(data)
	}
	return
}
