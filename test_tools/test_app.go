package test_tools

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/go-martini/martini"

    "github.com/hashtock/hashtock-go/conf"
    "github.com/hashtock/hashtock-go/models"
    "github.com/hashtock/hashtock-go/services"
    "github.com/hashtock/hashtock-go/webapp"
)

type TestApp struct {
    cfg     *conf.Config
    server  *httptest.Server
    test    *testing.T
    storage *models.MgoStorage

    AdminUser *models.Profile
    User      *models.Profile
    NoUser    *models.Profile
}

func NewTestApp(t *testing.T) *TestApp {
    var err error

    services.EnforceAuth = func(loginUrl string, exceptions ...string) martini.Handler {
        return func(req *http.Request, ctx martini.Context) { ctx.Next() }
    }

    app := new(TestApp)
    cfg := new(conf.Config)

    cfg.General = conf.GeneralConf{
        AppAddress:    "localhost:8123",
        DB:            "localhost",
        DBName:        fmt.Sprintf("test_db_%v", time.Now().UnixNano()),
        ServeAddr:     "localhost:8123",
        SessionKey:    "test_key",
        SessionSecret: "test_secter",
        Admin:         []string{"admin@a.com"},
    }

    app.storage, err = models.InitMongoStorage(cfg.General.DB, cfg.General.DBName)
    if err != nil {
        t.Fatal(err)
    }

    handler := webapp.Handlers(cfg, app.storage)
    app.server = httptest.NewServer(handler)
    app.test = t

    if app.AdminUser, err = app.CreateProfile("admin@a.com", true); err != nil {
        t.Fatal(err)
    }
    if app.User, err = app.CreateProfile("user@b.com", false); err != nil {
        t.Fatal(err)
    }

    return app
}

func (t *TestApp) Stop() {
    t.server.Close()
    t.storage.Collection("").Database.DropDatabase()
}
