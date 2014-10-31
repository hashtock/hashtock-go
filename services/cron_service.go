package services

import (
    "log"
    "net/http"

    "github.com/hashtock/hashtock-go/api"
    "github.com/hashtock/hashtock-go/http_utils"

    "github.com/hashtock/hashtock-go/models"
)

type CronService struct{}

func (c *CronService) Name() string {
    return "_cron"
}

func (c *CronService) EndPoints() (endpoints []*api.EndPoint) {
    bank_orders := api.NewEndPoint("/bank-orders/", "GET", "execute_bank_orders", ExecuteBankOrders)

    endpoints = []*api.EndPoint{
        bank_orders,
    }
    return
}

// Trigger execution of bank orders
// TODO(access): This endpoind should only by available to admin or cron
func ExecuteBankOrders(rw http.ResponseWriter, req *http.Request) {
    activeOrders, err := models.GetAllActiveBankOrders(req)
    if err != nil {
        http_utils.SimpleResponse(rw, err.Error(), http.StatusInternalServerError)
        return
    }

    //TODO(error): Handle errors somehow
    for _, order := range activeOrders {
        if err := models.ExecuteBankOrder(req, order); err != nil {
            log.Println(err)
        }
    }

    rw.WriteHeader(http.StatusOK)
}
