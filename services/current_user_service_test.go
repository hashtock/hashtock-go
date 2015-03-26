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

// TODO(security): It would be good to expand this to test ALL api urls
func TestUserHasToBeLoggedIn(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    expectedStatus := http.StatusForbidden
    expected := test_tools.Json{
        "code":  expectedStatus,
        "error": http.StatusText(expectedStatus),
    }

    rec := app.ExecuteJsonRequest("GET", "/api/user/", nil, app.NoUser)
    json_body := app.JsonResponceToStringMap(rec)

    json_body.Equal(t, expected)
    assert.Equal(t, expectedStatus, rec.Code)
}

func TestProfileForLoggedInUser(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    rec := app.ExecuteJsonRequest("GET", "/api/user/", nil, app.User)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "id":     app.User.UserID,
        "founds": models.StartingFounds,
    }

    assert.Equal(t, http.StatusOK, rec.Code)
    json_body.Equal(t, expected)
}
