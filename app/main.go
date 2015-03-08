package app

import (
    "net/http"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "github.com/hashtock/hashtock-go/services"
)

func init() {
    m := martini.Classic()
    m.Use(render.Renderer())
    m.Use(services.EnforceAuth("/auth/"))

    m.Group("/api", func(r martini.Router) {
        r.Group("/user", func(sr martini.Router) {
            sr.Get("/", services.CurrentProfile).Name("User:CurentUser")
        })

        r.Group("/portfolio", func(sr martini.Router) {
            sr.Get("/", services.Portfolio).Name("Portfolio:All")
            sr.Get("/:tag/", services.PortfolioTagInfo).Name("Portfolio:TagInfo")
        })

        r.Group("/tag", func(sr martini.Router) {
            sr.Get("/", services.ListOfAllHashTags).Name("Tag:Tags")
            sr.Post("/", services.NewHashTag).Name("Tag:newTag")
            sr.Get("/:tag/", services.TagInfo).Name("Tag:TagInfo")
            sr.Put("/:tag/", services.SetTagValue).Name("Tag:setTagValue")
            sr.Get("/:tag/values/", services.TagValues).Name("Tag:TagValues")
        })

        r.Group("/order", func(sr martini.Router) {
            sr.Get("/", services.ActiveOrders).Name("Order:Orders")
            sr.Post("/", services.NewOrder).Name("Order:NewOrder")
            sr.Get("/history/", services.CompletedOrder).Name("Order:CompletedOrders")
            sr.Get("/:uuid/", services.OrderDetails).Name("Order:OrderDetails")
            sr.Delete("/:uuid/", services.CancelOrder).Name("Order:CancelOrder")
        })

        r.Get("/", apiDefinition)
    })

    m.Group("/_cron", func(r martini.Router) {
        r.Get("/bank-orders/", services.ExecuteBankOrders).Name("Cron:executeBankOrders")
        r.Get("/tag-values/", services.FetchLatestTagValues).Name("Cron:fetchLatestTagValues")
    })

    m.Group("/auth", func(r martini.Router) {
        r.Get("/login/", login).Name("Auth:Login")
        r.Get("/logout/", logout).Name("Auth:Logout")
    })

    http.Handle("/", m)
}
