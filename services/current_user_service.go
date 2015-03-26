package services

import (
    "net/http"
    "strings"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/oauth2"
    "github.com/martini-contrib/render"
    "github.com/martini-contrib/sessions"

    "github.com/hashtock/hashtock-go/core"
    "github.com/hashtock/hashtock-go/models"
)

const (
    profileSessionKey = "hashtock_user"
)

func CurrentProfile(req *http.Request, r render.Render, tokens oauth2.Tokens, session sessions.Session) {
    profile, err := models.GetProfile(req)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusOK, profile)
}

func EnforceAuthFunc(loginUrl string, exceptions ...string) martini.Handler {
    return func(req *http.Request, rw http.ResponseWriter, ctx martini.Context, r render.Render, session sessions.Session, tokens oauth2.Tokens) {
        isException := false
        for _, prefix := range exceptions {
            if strings.HasPrefix(req.URL.Path, prefix) {
                isException = true
                break
            }
        }

        if !isException {
            if tokens.Expired() {
                r.Redirect(loginUrl)
                return
            }

            userId := ""
            if sessionValue := session.Get(profileSessionKey); sessionValue != nil {
                userId = sessionValue.(string)
            }

            profile, err := models.GetOrRegisterProfile(req, userId, tokens.Access())
            if err != nil {
                r.JSON(core.ErrToErrorer(err))
                return
            }

            session.Set(profileSessionKey, profile.UserID)
        }

        ctx.Next()
    }
}

var EnforceAuth = EnforceAuthFunc
