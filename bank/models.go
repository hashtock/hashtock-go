package bank

import (
    "appengine"
    "appengine/datastore"
)

type BankEntry struct {
    Hash   string
    Value  float64
    InBank float64
}

func Key(ctx appengine.Context, hash string) (key *datastore.Key) {
    return datastore.NewKey(ctx, "BankEntry", hash, 0, nil)
}

func (b *BankEntry) Key(ctx appengine.Context) (key *datastore.Key) {
    return Key(ctx, b.Hash)
}

func (b *BankEntry) Put(ctx appengine.Context) (err error) {
    key := b.Key(ctx)
    _, err = datastore.Put(ctx, key, b)
    return
}
