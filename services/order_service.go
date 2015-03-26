package services

import (
    "encoding/json"
    "net/http"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "github.com/hashtock/hashtock-go/core"
    "github.com/hashtock/hashtock-go/models"
)

func ActiveOrders(req *http.Request, r render.Render) {
    profile, err := models.GetProfile(req)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    orders, err := models.GetActiveUserOrders(req, profile)

    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusOK, orders)
}

func CompletedOrder(req *http.Request, r render.Render) {
    profile, err := models.GetProfile(req)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    queryValues := req.URL.Query()
    tag := queryValues.Get("tag")
    resolution := queryValues.Get("resolution")

    orders, err := models.GetCompletedUserOrders(req, profile, tag, resolution)

    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusOK, orders)
}

func NewOrder(req *http.Request, r render.Render) {
    profile, err := models.GetProfile(req)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    order := models.OrderBase{}
    decoder := json.NewDecoder(req.Body)
    if err := decoder.Decode(&order); err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    full_order, err := models.PlaceOrder(req, profile, order)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusCreated, full_order)
}

func OrderDetails(req *http.Request, params martini.Params, r render.Render) {
    uuid := params["uuid"]

    profile, err := models.GetProfile(req)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    order, err := models.GetOrder(req, profile, uuid)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusOK, order)
}

func CancelOrder(req *http.Request, params martini.Params, r render.Render) {
    uuid := params["uuid"]

    profile, err := models.GetProfile(req)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    if err := models.CancelOrder(req, profile, uuid); err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.Status(http.StatusNoContent)
}
