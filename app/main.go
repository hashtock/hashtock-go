package app

import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/rs/cors"

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
    cron_service := &services.CronService{}
    myapi := api.NewApi(app_routes, user_service, tag_service, order_service, cron_service)

    corsMiddleware := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowCredentials: true,
    })

    n := negroni.New(myapi.Middlewares()...)
    n.Use(corsMiddleware)
    n.UseHandler(r)

    http.Handle("/", n)
}
