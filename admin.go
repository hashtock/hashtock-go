package main

import (
    "errors"
    "fmt"
    "net/http"

    "appengine"
    "appengine/datastore"
    "appengine/user"
    "github.com/gorilla/mux"

    "bank"
)

func AddHashTag(ctx appengine.Context, hash string) (err error) {
    key := datastore.NewKey(ctx, "BankEntry", hash, 0, nil)

    entry := bank.BankEntry{
        Hash:   hash,
        Value:  1,
        InBank: 100,
    }

    if datastore.Get(ctx, key, &entry) != datastore.ErrNoSuchEntity {
        err = errors.New("Tag already exist")
        return
    }

    _, err = datastore.Put(ctx, key, &entry)
    return
}

func AddHashTagView(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    if !user.IsAdmin(ctx) {
        http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
        return
    }

    vars := mux.Vars(req)

    hash := vars["hash"]
    if err := AddHashTag(ctx, hash); err != nil {
        http.Error(rw, err.Error(), http.StatusNotFound)
    }

    rw.WriteHeader(http.StatusOK)
    fmt.Fprintf(rw, "OK: %#v", hash)
}
