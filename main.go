package main

import (
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

    r.HandleFunc("/admin/tag/add/{hash}", AddHashTagView)

    bank_routes := r.PathPrefix("/bank/").Subrouter()
    bank.AttachViews(bank_routes)

    profiles_routes := r.PathPrefix("/user/").Subrouter()
    profiles.AttachViews(profiles_routes)

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
