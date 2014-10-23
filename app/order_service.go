package main

import (
    "net/http"

    "github.com/gorilla/mux"

    "github.com/hashtock/hashtock-go/api"
    "github.com/hashtock/hashtock-go/http_utils"
    "github.com/hashtock/hashtock-go/models"
)

type OrderService struct{}

func (o *OrderService) Name() string {
    return "order"
}

func (o *OrderService) EndPoints() (endpoints []*api.EndPoint) {
    orders := api.NewEndPoint("/", "GET", "orders", ActiveOrder)
    new_order := api.NewEndPoint("/", "POST", "new_order", NewOrder)
    order_details := api.NewEndPoint("/{uuid}/", "GET", "order_details", OrderDetails)
    cancel_order := api.NewEndPoint("/{uuid}/", "DELETE", "cancel_order", CancelOrder)

    endpoints = []*api.EndPoint{
        orders,
        new_order,
        order_details,
        cancel_order,
    }
    return
}

func ActiveOrder(rw http.ResponseWriter, req *http.Request) {
    orders, err := models.GetActiveOrders(req)

    if err != nil {
        http_utils.SerializeErrorResponse(rw, req, err)
        return
    }

    http_utils.SerializeResponse(rw, req, orders, http.StatusOK)
}

func NewOrder(rw http.ResponseWriter, req *http.Request) {
    order := models.OrderBase{}

    if err := http_utils.DeSerializeRequest(*req, &order); err != nil {
        panic(err) // ToDo
        return
    }

    full_order, err := models.PlaceOrder(req, order)
    if err != nil {
        http_utils.SerializeErrorResponse(rw, req, err)
        return
    }

    http_utils.SerializeResponse(rw, req, full_order, http.StatusCreated)
}

func OrderDetails(rw http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    uuid := vars["uuid"]

    order, err := models.GetOrder(req, uuid)
    if err != nil {
        http_utils.SerializeErrorResponse(rw, req, err)
        return
    }

    http_utils.SerializeResponse(rw, req, order, http.StatusOK)
}

func CancelOrder(rw http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    uuid := vars["uuid"]

    if err := models.CancelOrder(req, uuid); err != nil {
        http_utils.SerializeErrorResponse(rw, req, err)
        return
    }

    http_utils.SerializeResponse(rw, req, nil, http.StatusNoContent)
}
