package main

import (
    "net/http"

    "github.com/hashtock/hashtock-go/api"
    "github.com/hashtock/hashtock-go/http_utils"
    "github.com/hashtock/hashtock-go/models"
)

type OrderService struct{}

func (o *OrderService) Name() string {
    return "order"
}

func (o *OrderService) EndPoints() (endpoints []*api.EndPoint) {
    orders := api.NewEndPoint("/", "GET", "orders", nil)
    new_order := api.NewEndPoint("/", "POST", "new_order", NewOrder)

    endpoints = []*api.EndPoint{
        orders,
        new_order,
    }
    return
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
