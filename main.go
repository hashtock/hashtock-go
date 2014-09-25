package main

import (
    "fmt"
    "net/http"

    "appengine"
    "appengine/user"
    "github.com/gorilla/mux"
)

func init() {
    r := mux.NewRouter()
    r.HandleFunc("/login", Login)
    r.HandleFunc("/transaction", TSell).Methods("POST")
    r.HandleFunc("/bank/sell/{hash}/{amount:\\d+}", BankSellHandler)
    r.HandleFunc("/bank/buy/{hash}/{amount}", BankBuyHandler)

    http.Handle("/", r)
}

func Login(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    u := user.Current(ctx)

    if u == nil {
        url, err := user.LoginURL(ctx, req.URL.String())
        if err != nil {
            http.Error(rw, err.Error(), http.StatusInternalServerError)
            return
        }
        rw.Header().Set("Location", url)
        rw.WriteHeader(http.StatusFound)
        return
    }
    fmt.Fprintf(rw, "Hello, %v!", u)
}

func BankSellHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    fmt.Fprintf(w, "Sell %v\n", vars["hash"])
    fmt.Fprintf(w, "Req %#v", r)
}

func BankBuyHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Buy")
}
