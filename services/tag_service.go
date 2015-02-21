package services

import (
    "encoding/json"
    "net/http"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "github.com/hashtock/hashtock-go/core"
    "github.com/hashtock/hashtock-go/models"
)

// List of all tags with bank values
func ListOfAllHashTags(req *http.Request, r render.Render) {
    tags, err := models.GetAllHashTags(req)

    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusOK, tags)
}

// Details about the hash tag
func TagInfo(req *http.Request, params martini.Params, r render.Render) {
    hash_tag_name := params["tag"]

    tag, err := models.GetHashTag(req, hash_tag_name)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusOK, tag)
}

// Add new tag (admin)
func NewHashTag(req *http.Request, r render.Render) {
    tag := models.HashTag{}

    decoder := json.NewDecoder(req.Body)
    if err := decoder.Decode(&tag); err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    new_tag, err := models.AddHashTag(req, tag)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusCreated, new_tag)
}

// Set tag value (admin)
func SetTagValue(req *http.Request, params martini.Params, r render.Render) {
    hash_tag_name := params["tag"]

    updated_tag := models.HashTag{}
    decoder := json.NewDecoder(req.Body)
    if err := decoder.Decode(&updated_tag); err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    if updated_tag.HashTag != "" && hash_tag_name != updated_tag.HashTag {
        err := core.NewBadRequestError("hashtag value has to be empty or correct")
        r.JSON(core.ErrToErrorer(err))
        return
    }

    tag, err := models.UpdateHashTagValue(req, hash_tag_name, updated_tag.Value)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusOK, tag)
}
