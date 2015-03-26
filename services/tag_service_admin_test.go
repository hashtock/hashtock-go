package services_test

import (
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/hashtock/hashtock-go/models"
    "github.com/hashtock/hashtock-go/test_tools"
)

func TestAdminAddingTag(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   10.5,
    }

    body := app.ToJsonBody(tag)
    req := app.NewJsonRequest("POST", "/api/tag/", body, app.AdminUser)
    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "hashtag": "TestTag",
        "value":   10.5,
        "in_bank": 100.0,
    }

    assert.Equal(t, http.StatusCreated, rec.Code)
    json_body.Equal(t, expected)

    tag.InBank = 100
    new_tag, err := models.GetHashTag(req, "TestTag")
    assert.NoError(t, err)
    assert.Equal(t, tag, *new_tag)
}

func TestAdminAddingTagWithInBankValue(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   10.5,
        InBank:  1.0,
    }

    body := app.ToJsonBody(tag)
    req := app.NewJsonRequest("POST", "/api/tag/", body, app.AdminUser)
    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "hashtag": "TestTag",
        "value":   10.5,
        "in_bank": 100.0,
    }

    assert.Equal(t, http.StatusCreated, rec.Code)
    json_body.Equal(t, expected)

    new_tag, err := models.GetHashTag(req, "TestTag")
    assert.NoError(t, err)
    assert.Equal(t, 100.0, new_tag.InBank)
}

func TestAdminAddingExistingTag(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   10.5,
        InBank:  100.0,
    }
    err := app.Put(tag)
    assert.NoError(t, err)

    body := app.ToJsonBody(tag)
    req := app.NewJsonRequest("POST", "/api/tag/", body, app.AdminUser)

    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "code":  400,
        "error": "Tag alread exists",
    }

    assert.Equal(t, http.StatusBadRequest, rec.Code)
    json_body.Equal(t, expected)
}

func TestAdminUpdateTagValue(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    existing_tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  50.0,
    }
    err := app.Put(existing_tag)
    assert.NoError(t, err)

    update_tag_value := models.HashTag{
        Value: 2,
    }

    body := app.ToJsonBody(update_tag_value)
    req := app.NewJsonRequest("PUT", "/api/tag/TestTag/", body, app.AdminUser)
    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "hashtag": "TestTag",
        "value":   2,
        "in_bank": 50.0,
    }

    assert.Equal(t, http.StatusOK, rec.Code)
    json_body.Equal(t, expected)
}

func TestAdminUpdateTagValueIgnoreBankValue(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    existing_tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  50.0,
    }
    err := app.Put(existing_tag)
    assert.NoError(t, err)

    update_tag_value := models.HashTag{
        Value:  2,
        InBank: 100.0, // This will be just ignored
    }

    body := app.ToJsonBody(update_tag_value)
    req := app.NewJsonRequest("PUT", "/api/tag/TestTag/", body, app.AdminUser)
    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "hashtag": "TestTag",
        "value":   2,
        "in_bank": 50.0,
    }

    assert.Equal(t, http.StatusOK, rec.Code)
    json_body.Equal(t, expected)
}

func TestAdminUpdateTagValueInvalidHashTag(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    existing_tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  50.0,
    }
    err := app.Put(existing_tag)
    assert.NoError(t, err)

    update_tag_value := models.HashTag{
        HashTag: "SomethingStupid",
        Value:   2,
    }

    body := app.ToJsonBody(update_tag_value)
    req := app.NewJsonRequest("PUT", "/api/tag/TestTag/", body, app.AdminUser)
    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "code":  400,
        "error": "hashtag value has to be empty or correct",
    }
    tag, _ := models.GetHashTag(req, "TestTag")

    assert.Equal(t, http.StatusBadRequest, rec.Code)
    json_body.Equal(t, expected)
    assert.Equal(t, existing_tag, *tag) // No change!!
}

func TestAdminUpdateTagInvalidValue(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    existing_tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  50.0,
    }
    err := app.Put(existing_tag)
    assert.NoError(t, err)

    update_tag_value := models.HashTag{
        Value: 0,
    }

    body := app.ToJsonBody(update_tag_value)
    req := app.NewJsonRequest("PUT", "/api/tag/TestTag/", body, app.AdminUser)
    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "code":  400,
        "error": "Value has to be positive",
    }
    tag, _ := models.GetHashTag(req, "TestTag")

    assert.Equal(t, http.StatusBadRequest, rec.Code)
    json_body.Equal(t, expected)
    assert.Equal(t, existing_tag, *tag) // No change!!
}

func TestRegularUserUpdateTagValue(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    existing_tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  50.0,
    }
    err := app.Put(existing_tag)
    assert.NoError(t, err)

    update_tag_value := models.HashTag{
        Value: 2,
    }

    body := app.ToJsonBody(update_tag_value)
    req := app.NewJsonRequest("PUT", "/api/tag/TestTag/", body, app.User)
    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "code":  http.StatusForbidden,
        "error": "Forbidden",
    }
    tag, _ := models.GetHashTag(req, "TestTag")

    assert.Equal(t, http.StatusForbidden, rec.Code)
    json_body.Equal(t, expected)
    assert.Equal(t, existing_tag, *tag) // No change!!
}

func TestRegularUserAddingTag(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   10.5,
    }

    body := app.ToJsonBody(tag)
    req := app.NewJsonRequest("POST", "/api/tag/", body, app.User)
    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "code":  http.StatusForbidden,
        "error": "Forbidden",
    }

    assert.Equal(t, http.StatusForbidden, rec.Code)
    json_body.Equal(t, expected)
}
