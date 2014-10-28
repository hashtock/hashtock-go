package main

import (
    "net/http"

    "github.com/gorilla/mux"

    "github.com/hashtock/hashtock-go/api"
    "github.com/hashtock/hashtock-go/http_utils"

    "github.com/hashtock/hashtock-go/models"
)

type HashTagService struct{}

func (h *HashTagService) Name() string {
    return "tag"
}

func (h *HashTagService) EndPoints() (endpoints []*api.EndPoint) {
    tags := api.NewEndPoint("/", "GET", "tags", ListOfAllHashTags)
    new_tag := api.NewEndPoint("/", "PUT", "new_tag", NewHashTag)
    tag_info := api.NewEndPoint("/{tag}/", "GET", "tag_info", TagInfo)
    set_tag_value := api.NewEndPoint("/{tag}/", "POST", "set_tag_value", SetTagValue)

    endpoints = []*api.EndPoint{
        tags,
        new_tag,
        tag_info,
        set_tag_value,
    }
    return
}

// List of all tags with bank values
func ListOfAllHashTags(rw http.ResponseWriter, req *http.Request) {
    tags, err := models.GetAllHashTags(req)

    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }

    http_utils.SerializeResponse(rw, req, tags, http.StatusOK)
}

// Details about the hash tag
func TagInfo(rw http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    hash_tag_name := vars["tag"]

    tag, err := models.GetHashTag(req, hash_tag_name)
    if err != nil {
        http_utils.SerializeErrorResponse(rw, req, err)
        return
    }

    http_utils.SerializeResponse(rw, req, tag, http.StatusOK)
}

// Add new tag (admin)
func NewHashTag(rw http.ResponseWriter, req *http.Request) {
    tag := models.HashTag{}
    if err := http_utils.DeSerializeRequest(*req, &tag); err != nil {
        http_utils.SerializeErrorResponse(rw, req, err)
        return
    }

    new_tag, err := models.AddHashTag(req, tag)
    if err != nil {
        http_utils.SerializeErrorResponse(rw, req, err)
        return
    }

    http_utils.SerializeResponse(rw, req, new_tag, http.StatusCreated)
}

// Set tag value (admin)
func SetTagValue(rw http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    hash_tag_name := vars["tag"]

    updated_tag := models.HashTag{}
    if err := http_utils.DeSerializeRequest(*req, &updated_tag); err != nil {
        http_utils.SerializeErrorResponse(rw, req, err)
        return
    }

    if updated_tag.HashTag != "" && hash_tag_name != updated_tag.HashTag {
        err := http_utils.NewBadRequestError("hashtag value has to be empty or correct")
        http_utils.SerializeErrorResponse(rw, req, err)
        return
    }

    tag, err := models.UpdateHashTagValue(req, hash_tag_name, updated_tag.Value)
    if err != nil {
        http_utils.SerializeErrorResponse(rw, req, err)
        return
    }

    http_utils.SerializeResponse(rw, req, tag, http.StatusOK)
}
