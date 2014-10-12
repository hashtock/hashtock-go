package main

import (
    "net/http"

    "github.com/gorilla/mux"

    "github.com/hashtock/hashtock-go/api"
)

func init() {
    r := mux.NewRouter()

    app_routes := r.PathPrefix("/api/").Subrouter()
    user_service := &CurrentUserService{}
    api.NewApi(app_routes, user_service)
    http.Handle("/", app_routes)
}
