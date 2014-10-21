package models

import (
    // "errors"
    "net/http"

    "appengine"
    "appengine/datastore"
    "code.google.com/p/go-uuid/uuid"
)

func orderKey(ctx appengine.Context, uuid string) (key *datastore.Key) {
    return datastore.NewKey(ctx, orderKind, uuid, 0, nil)
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
