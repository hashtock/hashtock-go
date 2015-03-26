package webapp

import (
    "net/http"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/oauth2"
    "github.com/martini-contrib/render"
    "github.com/martini-contrib/sessions"

    "github.com/hashtock/hashtock-go/conf"
    "github.com/hashtock/hashtock-go/models"
    "github.com/hashtock/hashtock-go/services"
)

// ToDo: Propagate this instance of storage down to services and models
func Handlers(cfg *conf.Config, storage *models.MgoStorage) http.Handler {
    authCfg := cfg.GAuthConfig()

    oauth2.PathLogin = "/auth/login/"
    oauth2.PathLogout = "/auth/logout/"
    oauth2.PathCallback = "/oauth2callback"

    cookieStore := sessions.NewCookieStore([]byte(cfg.General.SessionSecret))

    m := martini.Classic()
    m.Use(martini.Static("static", martini.StaticOptions{Prefix: "static"}))
    m.Use(martini.Static("static", martini.StaticOptions{IndexFile: "index.html"}))
    m.Use(sessions.Sessions(cfg.General.SessionKey, cookieStore))
    m.Use(oauth2.Google(&authCfg))
    m.Use(render.Renderer())
    m.Use(services.EnforceAuth(oauth2.PathLogin, "/auth/"))

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

    return m
}
