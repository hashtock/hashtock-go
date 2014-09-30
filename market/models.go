package market

import (
    "encoding/json"
    "errors"

    "appengine"
    "appengine/datastore"
)

type TransactionOrderBase struct {
    IsBuy   bool
    HashTag string
    Amount  float64
    Price   float64
}

type TransactionOrder struct {
    TransactionOrderBase

    UUID       string
    UserID     string
    Complete   bool
    ProfileKey *datastore.Key `json:"-"`
}

func Key(ctx appengine.Context, uuid string, profile_key *datastore.Key) (key *datastore.Key) {
    return datastore.NewKey(ctx, "TransactionOrder", uuid, 0, profile_key)
}

func (t *TransactionOrder) Key(ctx appengine.Context) (key *datastore.Key) {
    return Key(ctx, t.UUID, t.ProfileKey)
}

func (t *TransactionOrder) Put(ctx appengine.Context) (err error) {
    key := t.Key(ctx)
    _, err = datastore.Put(ctx, key, t)
    return
}

func (t *TransactionOrder) UnmarshalJSON(b []byte) (err error) {
    err = json.Unmarshal(b, &t.TransactionOrderBase)
    if err != nil {
        return
    }

    return t.IsValid()
}

func (t *TransactionOrder) IsValid() (err error) {
    if t.HashTag == "" {
        // ToDo: Add DB check
        err = errors.New("Invalid transaction: HashTag not specified")
        return
    }

    if t.Amount <= 0 || t.Amount > 100 {
        // ToDo: Check if user has enough shares
        err = errors.New("Invalid transaction: Amount outside range of (0, 100]")
        return
    }

    if t.Price <= 0 {
        // ToDo: If buying check if user has enough founds
        err = errors.New("Invalid transaction: Price lower than zero")
        return
    }

    return
}

func (t *TransactionOrder) CanDelete() error {
    if t.Complete {
        return errors.New("Transaction is complete. Could not be deleted")
    }
    return nil
}

func (t *TransactionOrder) Delete(ctx appengine.Context) (err error) {
    if err = t.CanDelete(); err != nil {
        return
    }

    err = datastore.Delete(ctx, t.Key(ctx))

    return
}
