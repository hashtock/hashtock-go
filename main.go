package main

import (
    "fmt"
    "net/http"

    "appengine"
    "appengine/user"
    "github.com/gorilla/mux"

    "bank"
    "profiles"
)

func init() {
    r := mux.NewRouter()
    r.HandleFunc("/login", Login)
    r.HandleFunc("/logout", Logout)

    r.HandleFunc("/bank/info/{hash}/", bank.HashTagInfoView)
    r.HandleFunc("/bank/info/", bank.HashTagInfoAllView)
    r.HandleFunc("/bank/sell_to_user/{hash}/{amount}/", bank.SellToUserView)
    r.HandleFunc("/bank/buy_from_user/{hash}/{amount}/", bank.BuyFromUserView)

    r.HandleFunc("/user/shares/", profiles.UserSharesView)

    r.HandleFunc("/admin/tag/add/{hash}", AddHashTagView)

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

    if err := profiles.CreateNewUserIfDoesNotExist(ctx, *u); err != nil {
        panic(err)
    }

    fmt.Fprintf(rw, "Hello, %v!", u)
}

func Logout(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)

    url, err := user.LogoutURL(ctx, "login")
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }
    rw.Header().Set("Location", url)
    rw.WriteHeader(http.StatusFound)
}
