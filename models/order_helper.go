package models

import (
    "errors"
    "fmt"
    "net/http"
    "time"

    "code.google.com/p/go-uuid/uuid"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    "github.com/hashtock/hashtock-go/core"
)

func newOrderSystem(req *http.Request, profile *Profile, base_order OrderBase) (order OrderSystem, err error) {
    var hashTag *HashTag

    if hashTag, err = GetHashTag(req, base_order.HashTag); err != nil {
        return
    }

    order = OrderSystem{
        UUID:       uuid.New(),
        UserID:     profile.UserID,
        Complete:   false,
        CreatedAt:  time.Now(),
        Resolution: PENDING,
        Value:      base_order.Quantity * hashTag.Value,
    }

    return
}

func PlaceOrder(req *http.Request, profile *Profile, base_order OrderBase) (order *Order, err error) {
    if err = baseOrderValid(req, base_order); err != nil {
        return
    }

    var system_order OrderSystem
    if system_order, err = newOrderSystem(req, profile, base_order); err != nil {
        return
    }

    order = &Order{
        OrderBase:   base_order,
        OrderSystem: system_order,
    }

    col := storage.Collection(OrderCollectionName)
    defer col.Database.Session.Close()

    err = col.Insert(order)

    return
}

func GetOrder(req *http.Request, profile *Profile, order_uuid string) (order *Order, err error) {
    col := storage.Collection(OrderCollectionName)
    defer col.Database.Session.Close()

    selector := bson.M{
        "user_id": profile.UserID,
        "uuid":    order_uuid,
    }

    err = col.Find(selector).One(&order)
    if err == mgo.ErrNotFound {
        not_found_msg := fmt.Sprintf("Order %#v not found", order_uuid)
        err = core.NewNotFoundError(not_found_msg)
    }
    return
}

func CancelOrder(req *http.Request, profile *Profile, order_uuid string) (err error) {
    order, err := GetOrder(req, profile, order_uuid)
    if err != nil {
        return
    }

    if order.Complete {
        err = core.NewBadRequestError("Order can not be cancelled any more")
        return
    }

    err = orderDelete(order)

    return
}

func GetActiveUserOrders(req *http.Request, profile *Profile) (orders []Order, err error) {
    col := storage.Collection(OrderCollectionName)
    defer col.Database.Session.Close()

    selector := bson.M{
        "user_id":  profile.UserID,
        "complete": false,
    }
    err = col.Find(selector).Sort("-created_at").All(&orders)
    return
}

func GetCompletedUserOrders(req *http.Request, profile *Profile, tag string, resolution string) (orders []Order, err error) {
    col := storage.Collection(OrderCollectionName)
    defer col.Database.Session.Close()

    selector := bson.M{
        "user_id":  profile.UserID,
        "complete": true,
    }

    if tag != "" {
        selector["hashtag"] = tag
    }

    if resolution != "" {
        selector["resolution"] = resolution
    }

    err = col.Find(selector).Sort("-created_at").All(&orders)
    return
}

func GetAllActiveBankOrders(req *http.Request) (orders []Order, err error) {
    col := storage.Collection(OrderCollectionName)
    defer col.Database.Session.Close()

    selector := bson.M{
        "complete":   false,
        "bank_order": true,
    }
    err = col.Find(selector).Sort("-created_at").All(&orders)
    return
}

func ExecuteBankOrder(req *http.Request, order Order) (err error) {
    var (
        profile  *Profile
        hashTag  *HashTag
        tagShare *TagShare
    )

    // It's time to blow up if asked to execute non bank order here
    if !order.BankOrder {
        panic(errors.New("execution of non bank order"))
    }

    if hashTag, err = GetHashTag(req, order.HashTag); err != nil {
        markOrderAsComplete(order, ERROR, "")
        return
    }

    if profile, err = getProfileForUserId(req, order.UserID); err != nil {
        markOrderAsComplete(order, ERROR, "")
        return
    }

    if tagShare, err = getOrCreateTagShare(req, profile, order.HashTag); err != nil {
        if err != mgo.ErrNotFound {
            markOrderAsComplete(order, ERROR, "")
            return
        }
        err = nil
    }

    if order.Action == actionBuy {
        if profile.Founds < order.Value {
            markOrderAsComplete(order, FAILURE, "Not enough founds")
            msg := fmt.Sprintf("User %v does not have enough founds to complete %v", profile, order)
            return core.NewBadRequestError(msg)
        }

        if hashTag.InBank < order.Quantity {
            markOrderAsComplete(order, FAILURE, "Not enough shares in bank")
            msg := fmt.Sprintf("Bank does not have enough shares to complete %v", order)
            return core.NewBadRequestError(msg)
        }

        profileUpdateFounds(profile, -order.Value)
        tagShareUpdateQuantity(tagShare, order.Quantity)
        hashTagUpdateInBank(hashTag, -order.Quantity)
    } else {
        if tagShare.Quantity < order.Quantity {
            markOrderAsComplete(order, FAILURE, "Not enough shares in users possession")
            msg := fmt.Sprintf("User %v does not have enough shares (%v) to complete %v - %#v", profile.UserID, tagShare.Quantity, order.UUID, order.OrderBase)
            return core.NewBadRequestError(msg)
        }

        if tagShare.Quantity < minShareStep/100.0 {
            tagShareDelete(tagShare)
        } else {
            tagShareUpdateQuantity(tagShare, -order.Quantity)
        }

        profileUpdateFounds(profile, order.Value)
        hashTagUpdateInBank(hashTag, order.Quantity)
    }

    markOrderAsComplete(order, SUCCESS, "")

    return
}
