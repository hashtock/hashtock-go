// Service wide tests
package services_test

import (
    "net/http"
    "testing"

    _ "github.com/hashtock/hashtock-go/app" // Imported to initialize app
    "github.com/hashtock/hashtock-go/gaetestsuite"
)

type ServicesTestSuite struct {
    gaetestsuite.GAETestSuite
}

func (s *ServicesTestSuite) TestApiHasAllEndpoints() {
    rec := s.ExecuteJsonRequest("GET", "/api/", nil, s.User)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "Auth:Login":            "/auth/login/",
        "Auth:Logout":           "/auth/logout/",
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

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

/* Kickoff Test Suite */

func TestServicesTestSuite(t *testing.T) {
    gaetestsuite.Run(t, new(ServicesTestSuite))
}
