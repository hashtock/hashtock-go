package services

import (
    "net/http"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "github.com/hashtock/hashtock-go/core"
    "github.com/hashtock/hashtock-go/models"
)

func CurrentProfile(req *http.Request, r render.Render) {
    profile, _ := models.GetProfile(req)

    r.JSON(http.StatusOK, profile)
}

func Shares(req *http.Request, r render.Render) {
    profile, _ := models.GetProfile(req)
    shares, _ := models.GetProfileShares(req, profile)

    r.JSON(http.StatusOK, shares)
}

func EnforceAuth(req *http.Request, c martini.Context, r render.Render) {
    if _, err := models.GetProfile(req); err != nil && req.Header.Get("X-AppEngine-Cron") == "" {
        hErr := core.NewForbiddenError()
        r.JSON(hErr.ErrCode(), hErr)
        return
    }

    c.Next()
}
