package main

import (
    "github.com/hashtock/hashtock-go/api"
)

type CurrentUserService struct {
}

func (c *CurrentUserService) Name() string {
    return "user"
}

func (c *CurrentUserService) EndPoints() (endpoints []*api.EndPoint) {
    user := api.NewEndPoint("/", "GET", "user", nil)           // High level user details
    tags := api.NewEndPoint("/tags/", "GET", "user_tags", nil) // List of users shares of tags

    endpoints = []*api.EndPoint{
        user,
        tags,
    }
    return
}
