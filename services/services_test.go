// Service wide tests
package services_test

import (
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/hashtock/hashtock-go/test_tools"
)

func TestApiHasAllEndpoints(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    rec := app.ExecuteJsonRequest("GET", "/api/", nil, app.User)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        // "Auth:Login":            "/auth/login/",
        // "Auth:Logout":           "/auth/logout/",
        "Order:CancelOrder":     "/api/order/:uuid/",
        "Order:CompletedOrders": "/api/order/history/",
        "Order:NewOrder":        "/api/order/",
        "Order:OrderDetails":    "/api/order/:uuid/",
        "Order:Orders":          "/api/order/",
        "Portfolio:All":         "/api/portfolio/",
        "Portfolio:TagInfo":     "/api/portfolio/:tag/",
        "Tag:TagInfo":           "/api/tag/:tag/",
        "Tag:Tags":              "/api/tag/",
        "Tag:TagValues":         "/api/tag/:tag/values/",
        "User:CurentUser":       "/api/user/",
    }

    assert.Equal(t, http.StatusOK, rec.Code)
    json_body.Equal(t, expected)
}
