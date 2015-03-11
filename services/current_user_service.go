package services

import (
    "net/http"
    "strings"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "github.com/hashtock/hashtock-go/core"
    "github.com/hashtock/hashtock-go/models"
)

func CurrentProfile(req *http.Request, r render.Render) {
    profile, _ := models.GetProfile(req)

    r.JSON(http.StatusOK, profile)
}

func EnforceAuth(exceptions ...string) martini.Handler {
    return func(req *http.Request, c martini.Context, r render.Render) {
        isException := false
        for _, prefix := range exceptions {
            if strings.HasPrefix(req.URL.Path, prefix) {
                isException = true
                break
            }
        }

        _, err := models.GetProfile(req)
        if !isException && req.Header.Get("X-AppEngine-Cron") == "" && err != nil {
            hErr := core.NewForbiddenError()
            r.JSON(hErr.ErrCode(), hErr)
            return
        }
        c.Next()
    }
}
