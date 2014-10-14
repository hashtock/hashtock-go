package main

import (
    "net/http"

    "github.com/gorilla/mux"

    "github.com/codegangsta/negroni"
    "github.com/hashtock/hashtock-go/api"
)

func init() {
    r := mux.NewRouter()
    app_routes := r.PathPrefix("/api/").Subrouter()

    user_service := &CurrentUserService{}
    // myapi := api.NewApi(app_routes, user_service)
    api.NewApi(app_routes, user_service)

    n := negroni.New(user_service)
    n.UseHandler(r)

    http.Handle("/", n)
    // http.Handle("/", r)
}
