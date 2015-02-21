package services

import (
    "log"
    "net/http"

    "github.com/martini-contrib/render"

    "github.com/hashtock/hashtock-go/models"
)

// Trigger execution of bank orders
// TODO(access): This endpoind should only by available to admin or cron
func ExecuteBankOrders(req *http.Request, r render.Render) {
    activeOrders, err := models.GetAllActiveBankOrders(req)
    if err != nil {
        r.Error(http.StatusInternalServerError)
        return
    }

    //TODO(error): Handle errors somehow
    for _, order := range activeOrders {
        if err := models.ExecuteBankOrder(req, order); err != nil {
            log.Println(err)
        }
    }

    r.Status(http.StatusOK)
}
