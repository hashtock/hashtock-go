package market

import (
    "appengine"
    "appengine/datastore"
    "code.google.com/p/go-uuid/uuid"

    "profiles"
)

func NewOrder() (order TransactionOrder) {
    order.UUID = uuid.New()
    return
}

func NewOrderForContext(ctx appengine.Context) (order TransactionOrder) {
    current_user, err := profiles.CurrentProfile(ctx)
    if err != nil {
        panic(err)
        return
    }
    order = NewOrder()
    order.UserID = current_user.UserID
    order.ProfileKey = current_user.Key(ctx)
    return
}

func GetAll(ctx appengine.Context) (entries []TransactionOrder, err error) {
    profile_key := profiles.CurrentProfileKey(ctx)

    q := datastore.NewQuery("TransactionOrder").Ancestor(profile_key)
    _, err = q.GetAll(ctx, &entries)
    return
}

func GetById(ctx appengine.Context, uuid string) (order TransactionOrder, err error) {
    current_user, err := profiles.CurrentProfile(ctx)
    if err != nil {
        return
    }
    key := Key(ctx, uuid, current_user.Key(ctx))
    if err = datastore.Get(ctx, key, &order); err != nil {
        return
    }

    return
}

func DeleteById(ctx appengine.Context, uuid string) (err error) {
    order, err := GetById(ctx, uuid)
    if err != nil {
        return
    }

    err = order.Delete(ctx)
    return
}
