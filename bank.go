package main

import (
    "errors"
    "fmt"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

type Bank struct {
}

var bank Bank

func (b *Bank) sell(hash string, amount float64) (err error) {
    if 0 >= amount || amount > 100 {
        err = errors.New("Amount outside boundries")
        return
    }

    return
}

func SellView(rw http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)

    amount, _ := strconv.ParseFloat(vars["amount"], 64)
    hash := vars["hash"]

    err := bank.sell(hash, amount)

    if err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Fprintf(rw, "Sell %v\n", vars["hash"])
}
