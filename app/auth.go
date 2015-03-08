package app

import (
    "log"
    "net/http"

    "appengine"
    "appengine/user"
)

func login(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    u := user.Current(ctx)

    if u == nil {
        values := req.URL.Query()
        redir := values.Get("continue")
        if redir == "" {
            redir = "/"
        }
        log.Println("URL: ", redir)

        url, err := user.LoginURL(ctx, redir)
        if err != nil {
            http.Error(rw, err.Error(), http.StatusInternalServerError)
            return
        }
        rw.Header().Set("Location", url)
    } else {
        rw.Header().Set("Location", "/")
    }

    rw.WriteHeader(http.StatusFound)
}

func logout(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)

    url, err := user.LogoutURL(ctx, "/")
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }
    rw.Header().Set("Location", url)
    rw.WriteHeader(http.StatusFound)
}
