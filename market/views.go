package market

import (
    "net/http"

    "appengine"
    "github.com/gorilla/mux"

    "http_utils"
)

func AcceptNewOrder(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    order := NewOrderForContext(ctx)

    if err := http_utils.DeSerializeRequest(*req, &order); err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

    if err := order.Put(ctx); err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

    http_utils.SerializeResponse(rw, req, order, http.StatusCreated)
}

func ViewOrder(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    vars := mux.Vars(req)

    uuid := vars["uuid"]

    order, err := GetById(ctx, uuid)
    if err != nil {
        http_utils.SerializeResponse(rw, req, err.Error(), http.StatusNotFound)
        return
    }

    http_utils.SerializeResponse(rw, req, order, http.StatusOK)
}

func DeleteOrder(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    vars := mux.Vars(req)

    uuid := vars["uuid"]

    if err := DeleteById(ctx, uuid); err != nil {
        http_utils.SerializeResponse(rw, req, err.Error(), http.StatusNotFound)
        return
    }

    http_utils.SerializeResponse(rw, req, "Deleted", http.StatusOK)
}

func ViewOrders(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)

    orders, err := GetAll(ctx)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusNotFound)
        return
    }

    http_utils.SerializeResponse(rw, req, orders, http.StatusOK)
}
