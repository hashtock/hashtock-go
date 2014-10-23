package models

import (
    "fmt"
    "net/http"

    "appengine"
    "appengine/datastore"
    "code.google.com/p/go-uuid/uuid"

    "github.com/hashtock/hashtock-go/http_utils"
)

func orderKey(ctx appengine.Context, order_uuid string) (key *datastore.Key) {
    return datastore.NewKey(ctx, orderKind, order_uuid, 0, nil)
}

func newOrderSystem(req *http.Request) (order OrderSystem, err error) {
    var profile *Profile

    if profile, err = GetProfile(req); err != nil {
        return
    }

    order = OrderSystem{
        UUID:     uuid.New(),
        UserID:   profile.UserID,
        Complete: false,
    }

    return
}

func PlaceOrder(req *http.Request, base_order OrderBase) (order *Order, err error) {
    if err = base_order.IsValid(req); err != nil {
        return
    }

    var system_order OrderSystem
    if system_order, err = newOrderSystem(req); err != nil {
        return
    }

    order = &Order{
        OrderBase:   base_order,
        OrderSystem: system_order,
    }

    order.Put(req)

    return
}

func GetOrder(req *http.Request, order_uuid string) (order *Order, err error) {
    ctx := appengine.NewContext(req)

    order = new(Order)
    key := orderKey(ctx, order_uuid)
    err = datastore.Get(ctx, key, order)

    not_found_msg := fmt.Sprintf("Order %#v not found", order_uuid)

    if err == datastore.ErrNoSuchEntity {
        err = http_utils.NewNotFoundError(not_found_msg)
    } else if err != nil {
        err = http_utils.NewInternalServerError(err.Error())
    }

    var ok bool
    ok, err = order.canAccess(req)
    if err != nil {
        err = http_utils.NewInternalServerError(err.Error())
    }

    if !ok {
        err = http_utils.NewNotFoundError(not_found_msg)
        order = nil
    }

    return
}

func CancelOrder(req *http.Request, order_uuid string) (err error) {
    order, err := GetOrder(req, order_uuid)
    if err != nil {
        return
    }

    if !order.isCancellable() {
        err = http_utils.NewBadRequestError("Order can not be canceled any more")
        return
    }

    order.Delete(req)

    return
}

func allUserOrders(req *http.Request) (query *datastore.Query, err error) {
    var profile *Profile

    if profile, err = GetProfile(req); err != nil {
        return
    }

    query = datastore.NewQuery(orderKind).Filter("UserID =", profile.UserID)

    return
}

func GetActiveOrders(req *http.Request) (orders []Order, err error) {
    var query *datastore.Query
    query, err = allUserOrders(req)

    ctx := appengine.NewContext(req)
    _, err = query.Filter("Complete =", false).GetAll(ctx, &orders)

    return
}
