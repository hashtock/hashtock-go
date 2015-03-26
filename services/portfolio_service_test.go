// Current user service
// Run as part of service test suite
package services_test

import (
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/hashtock/hashtock-go/models"
    "github.com/hashtock/hashtock-go/test_tools"
)

func TestGetUsersPortfolioOfTags(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/portfolio/", nil, app.User)

    app.Put(
        models.TagShare{HashTag: "Tag1", Quantity: 10.5, UserID: app.User.UserID},
        models.TagShare{HashTag: "bTag", Quantity: 0.20, UserID: app.User.UserID},
        models.TagShare{HashTag: "Tag3", Quantity: 1.20, UserID: app.User.UserID},
        models.TagShare{HashTag: "aTag", Quantity: 0.20, UserID: app.User.UserID},
        models.TagShare{HashTag: "Tag1", Quantity: 1.00, UserID: "OtherID"},
    )

    rec := app.Do(req)
    json_body := app.JsonResponceToListOfStringMap(rec)

    // Order matters
    expected := test_tools.JsonList{
        test_tools.Json{
            "hashtag":  "Tag1",
            "quantity": 10.5,
        },
        test_tools.Json{
            "hashtag":  "Tag3",
            "quantity": 1.2,
        },
        test_tools.Json{
            "hashtag":  "aTag",
            "quantity": 0.2,
        },
        test_tools.Json{
            "hashtag":  "bTag",
            "quantity": 0.2,
        },
    }

    assert.Equal(t, http.StatusOK, rec.Code)
    json_body.Equal(t, expected)
}

func TestGetTagShareFromPortfolio(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/portfolio/Tag1/", nil, app.User)

    app.Put(
        models.TagShare{HashTag: "Tag1", Quantity: 10.5, UserID: app.User.UserID},
        models.TagShare{HashTag: "Tag1", Quantity: 1.00, UserID: "OtherID"},
    )

    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "hashtag":  "Tag1",
        "quantity": 10.5,
    }

    assert.Equal(t, http.StatusOK, rec.Code)
    json_body.Equal(t, expected)
}

func TestGetTagShareFromPortfolioWrong(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/portfolio/FAKE/", nil, app.User)

    app.Put(
        models.TagShare{HashTag: "Tag1", Quantity: 10.5, UserID: app.User.UserID},
        models.TagShare{HashTag: "Tag1", Quantity: 1.00, UserID: "OtherID"},
    )

    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "error": "Not Found",
        "code":  http.StatusNotFound,
    }

    assert.Equal(t, http.StatusNotFound, rec.Code)
    json_body.Equal(t, expected)
}
