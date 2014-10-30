package services

import (
    "net/http"

    "github.com/hashtock/hashtock-go/api"
    "github.com/hashtock/hashtock-go/http_utils"

    "github.com/hashtock/hashtock-go/models"
)

type CurrentUserService struct{}

func (c *CurrentUserService) Name() string {
    return "user"
}

func (c *CurrentUserService) EndPoints() (endpoints []*api.EndPoint) {
    user := api.NewEndPoint("/", "GET", "user", c.Profile)          // High level user details
    tags := api.NewEndPoint("/tags/", "GET", "user_tags", c.Shares) // List of users shares of tags

    endpoints = []*api.EndPoint{
        user,
        tags,
    }
    return
}

func (c *CurrentUserService) Profile(rw http.ResponseWriter, req *http.Request) {
    profile, _ := models.GetProfile(req)

    http_utils.SerializeResponse(rw, req, profile, http.StatusOK)
}

func (c *CurrentUserService) Shares(rw http.ResponseWriter, req *http.Request) {
    profile, _ := models.GetProfile(req)
    shares, _ := models.GetProfileShares(req, profile)

    http_utils.SerializeResponse(rw, req, shares, http.StatusOK)
}

func (c *CurrentUserService) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
    if _, err := models.GetProfile(req); err != nil {
        http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
        return
    }

    next(rw, req)
}
