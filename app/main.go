package app

import (
    "net/http"

    "github.com/gorilla/mux"

    "github.com/codegangsta/negroni"
    "github.com/hashtock/hashtock-go/api"
    "github.com/hashtock/hashtock-go/services"
)

func init() {
    r := mux.NewRouter()
    app_routes := r.PathPrefix("/api/").Subrouter()

    user_service := &services.CurrentUserService{}
    tag_service := &services.HashTagService{}
    order_service := &services.OrderService{}
    myapi := api.NewApi(app_routes, user_service, tag_service, order_service)

    n := negroni.New(myapi.Middlewares()...)
    n.UseHandler(r)

    http.Handle("/", n)
}
