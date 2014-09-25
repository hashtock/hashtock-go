package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"

    "appengine"
    "appengine/user"
)

type TransactionQuery struct {
    User   string  `json:"-"`
    Hash   string  `json:"hash"`
    Amount float64 `json:"amount"`
    Price  float64 `json:"price"`
    IsSell bool    `json:"issell"`
}

func (t *TransactionQuery) BankTransaction() bool {
    if t.Price == 0 {
        return true
    }

    return false
}

func (t *TransactionQuery) IsValid() (ok bool, err error) {
    if 0 >= t.Amount || t.Amount > 100 {
        err = errors.New("Amount outside boundaries")
        ok = false
        return
    }

    if t.Price < 0 {
        err = errors.New("Price has to be positive")
        ok = false
        return
    }

    // ToDo
    // - Check is this a valid tag
    // - If sell does user have enough

    ok = true
    return
}

func TSell(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    u := user.Current(ctx)

    t := TransactionQuery{
        User: u.String(),
    }

    decoder := json.NewDecoder(req.Body)

    if err := decoder.Decode(&t); err != nil {
        panic(err)
    }

    if _, err := t.IsValid(); err != nil {
        panic(err)
    }

    fmt.Fprintf(rw, "POST:\n %#v\n", t)
    fmt.Fprintf(rw, "Bank: %#v\n", t.BankTransaction())
}
