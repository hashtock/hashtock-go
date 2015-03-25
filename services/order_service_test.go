// Order service
// Run as part of service test suite
package services_test

import (
    "net/http"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"

    "github.com/hashtock/hashtock-go/models"
    "github.com/hashtock/hashtock-go/test_tools"
)

func TestPlaceTransactionOrderWithBank(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    order := models.OrderBase{
        Action:    "buy",
        BankOrder: true,
        HashTag:   "Tag1",
        Quantity:  1.00,
    }

    body := app.ToJsonBody(order)
    req := app.NewJsonRequest("POST", "/api/order/", body, app.User)

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    app.Put(tag)

    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)
    assert.Equal(t, http.StatusCreated, rec.Code)

    uuid := json_body["uuid"]
    created := json_body["created_at"]
    expected := test_tools.Json{
        "action":      "buy",
        "hashtag":     "Tag1",
        "quantity":    1.00,
        "user_id":     app.User.UserID,
        "bank_order":  true,
        "complete":    false,
        "value":       1.00,
        "uuid":        uuid.(string),
        "created_at":  created,
        "executed_at": "0001-01-01T00:00:00Z", // Not executed yet
        "resolution":  "",
        "notes":       "",
    }
    order_in_db, err := models.GetOrder(req, app.User, uuid.(string))
    assert.NoError(t, err)
    assert.NotNil(t, order_in_db)
    json_body.Equal(t, expected)
}

func TestPlaceInvalidTransactionOrderWithBank(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    order := models.OrderBase{
        Action:    "freebe",
        BankOrder: true,
        HashTag:   "",
        Quantity:  101.00,
    }

    body := app.ToJsonBody(order)
    req := app.NewJsonRequest("POST", "/api/order/", body, app.User)

    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "code":  400,
        "error": "Incorrect fields: action, hashtag, quantity",
    }

    assert.Equal(t, http.StatusBadRequest, rec.Code)
    json_body.Equal(t, expected)
}

func TestGetOrderDetails(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/order/FAKE-UUID/", nil, app.User)

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    app.Put(tag)

    order := models.Order{
        OrderBase: models.OrderBase{
            Action:    "buy",
            BankOrder: true,
            HashTag:   "Tag1",
            Quantity:  1.00,
        },
        OrderSystem: models.OrderSystem{
            UUID:       "FAKE-UUID",
            UserID:     app.User.UserID,
            Complete:   false,
            Value:      1,
            CreatedAt:  time.Date(2015, time.February, 04, 19, 30, 0, 0, time.UTC),
            ExecutedAt: time.Time{},
            Resolution: models.PENDING,
            Notes:      "Some note",
        },
    }
    time.Local = time.UTC
    app.Put(order)

    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "action":      "buy",
        "hashtag":     "Tag1",
        "quantity":    1.00,
        "user_id":     app.User.UserID,
        "bank_order":  true,
        "complete":    false,
        "value":       1.00,
        "uuid":        "FAKE-UUID",
        "created_at":  "2015-02-04T19:30:00Z",
        "executed_at": "0001-01-01T00:00:00Z", // Not executed yet
        "resolution":  "",
        "notes":       "Some note",
    }

    assert.Equal(t, http.StatusOK, rec.Code)
    json_body.Equal(t, expected)
}

func TestGetOrderDifferentUser(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/order/FAKE-UUID/", nil, app.User)

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    app.Put(tag)

    order := models.Order{
        OrderBase: models.OrderBase{
            Action:    "buy",
            BankOrder: true,
            HashTag:   "Tag1",
            Quantity:  1.00,
        },
        OrderSystem: models.OrderSystem{
            UUID:     "FAKE-UUID",
            UserID:   "SOME USER",
            Complete: false,
        },
    }
    app.Put(order)

    rec := app.Do(req)

    assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCancelOrder(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("DELETE", "/api/order/FAKE-UUID/", nil, app.User)

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    app.Put(tag)

    order := models.Order{
        OrderBase: models.OrderBase{
            Action:    "buy",
            BankOrder: true,
            HashTag:   "Tag1",
            Quantity:  1.00,
        },
        OrderSystem: models.OrderSystem{
            UUID:     "FAKE-UUID",
            UserID:   app.User.UserID,
            Complete: false,
        },
    }
    app.Put(order)

    rec := app.Do(req)

    assert.Equal(t, http.StatusNoContent, rec.Code)
    assert.Equal(t, 0, rec.Body.Len(), rec.Body.String())

    cancelled_order, _ := models.GetOrder(req, app.User, "FAKE-UUID")
    assert.Nil(t, cancelled_order)
}

func TestCancelCompetedOrder(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("DELETE", "/api/order/FAKE-UUID/", nil, app.User)

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    app.Put(tag)

    order := models.Order{
        OrderBase: models.OrderBase{
            Action:    "buy",
            BankOrder: true,
            HashTag:   "Tag1",
            Quantity:  1.00,
        },
        OrderSystem: models.OrderSystem{
            UUID:     "FAKE-UUID",
            UserID:   app.User.UserID,
            Complete: true,
        },
    }
    app.Put(order)

    rec := app.Do(req)

    assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCurrentOrders(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/order/", nil, app.User)
    time.Local = time.UTC

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    app.Put(tag)

    base_order := models.OrderBase{
        Action:    "buy",
        BankOrder: true,
        HashTag:   "Tag1",
        Quantity:  1.00,
    }

    order_1 := models.Order{
        OrderBase: base_order,
        OrderSystem: models.OrderSystem{
            UUID:       "pending-1",
            UserID:     app.User.UserID,
            Complete:   false,
            Value:      1,
            CreatedAt:  time.Date(2015, time.February, 04, 19, 30, 0, 0, time.UTC),
            ExecutedAt: time.Time{},
        },
    }
    app.Put(order_1)

    order_1.UUID = "pending-2"
    order_1.CreatedAt = time.Date(2015, time.February, 04, 19, 31, 0, 0, time.UTC)
    order_1.ExecutedAt = time.Time{}
    app.Put(order_1)

    order_2 := models.Order{
        OrderBase: base_order,
        OrderSystem: models.OrderSystem{
            UUID:     "complete-1",
            UserID:   app.User.UserID,
            Complete: true,
        },
    }
    app.Put(order_2)

    order_3 := models.Order{
        OrderBase: base_order,
        OrderSystem: models.OrderSystem{
            UUID:     "pending-3",
            UserID:   "some user",
            Complete: false,
        },
    }
    app.Put(order_3)

    rec := app.Do(req)
    json_body := app.JsonResponceToListOfStringMap(rec)

    expected := test_tools.JsonList{
        test_tools.Json{
            "action":      "buy",
            "hashtag":     "Tag1",
            "quantity":    1.00,
            "user_id":     app.User.UserID,
            "bank_order":  true,
            "complete":    false,
            "uuid":        "pending-2",
            "value":       1.00,
            "created_at":  "2015-02-04T19:31:00Z",
            "executed_at": "0001-01-01T00:00:00Z", // Not executed yet
            "resolution":  "",
            "notes":       "",
        },
        test_tools.Json{
            "action":      "buy",
            "hashtag":     "Tag1",
            "quantity":    1.00,
            "user_id":     app.User.UserID,
            "bank_order":  true,
            "complete":    false,
            "uuid":        "pending-1",
            "value":       1.00,
            "created_at":  "2015-02-04T19:30:00Z",
            "executed_at": "0001-01-01T00:00:00Z", // Not executed yet
            "resolution":  "",
            "notes":       "",
        },
    }

    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Len(t, json_body, 2)
    json_body.Equal(t, expected)
}

func TestHistoricOrders(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/order/history/", nil, app.User)

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    app.Put(tag)

    base_order := models.OrderBase{
        Action:    "buy",
        BankOrder: true,
        HashTag:   "Tag1",
        Quantity:  1.00,
    }

    order_1 := models.Order{
        OrderBase: base_order,
        OrderSystem: models.OrderSystem{
            UUID:       "complete-1",
            UserID:     app.User.UserID,
            Complete:   true,
            Value:      1.0,
            CreatedAt:  time.Date(2015, time.February, 04, 19, 00, 0, 0, time.UTC),
            ExecutedAt: time.Date(2015, time.February, 04, 19, 15, 0, 0, time.UTC),
            Resolution: models.SUCCESS,
        },
    }
    app.Put(order_1)

    order_1.UUID = "complete-2"
    order_1.CreatedAt = time.Date(2015, time.February, 04, 19, 30, 0, 0, time.UTC)
    order_1.ExecutedAt = time.Date(2015, time.February, 04, 19, 45, 0, 0, time.UTC)
    order_1.Resolution = models.FAILURE
    order_1.Notes = "some reason"
    app.Put(order_1)

    order_2 := models.Order{
        OrderBase: base_order,
        OrderSystem: models.OrderSystem{
            UUID:       "pending-1",
            UserID:     app.User.UserID,
            Complete:   false,
            Value:      1.0,
            CreatedAt:  time.Date(2015, time.February, 04, 19, 00, 0, 0, time.UTC),
            ExecutedAt: time.Time{},
        },
    }
    app.Put(order_2)

    order_3 := models.Order{
        OrderBase: base_order,
        OrderSystem: models.OrderSystem{
            UUID:       "complete-3",
            UserID:     "some user",
            Complete:   true,
            Value:      1.0,
            CreatedAt:  time.Date(2015, time.February, 04, 19, 31, 0, 0, time.UTC),
            ExecutedAt: time.Date(2015, time.February, 04, 19, 46, 0, 0, time.UTC),
        },
    }
    app.Put(order_3)

    rec := app.Do(req)
    json_body := app.JsonResponceToListOfStringMap(rec)

    expected := test_tools.JsonList{
        test_tools.Json{
            "action":      "buy",
            "hashtag":     "Tag1",
            "quantity":    1.00,
            "user_id":     app.User.UserID,
            "bank_order":  true,
            "complete":    true,
            "value":       1.0,
            "uuid":        "complete-2",
            "created_at":  "2015-02-04T19:30:00Z",
            "executed_at": "2015-02-04T19:45:00Z",
            "resolution":  "failure",
            "notes":       "some reason",
        },
        test_tools.Json{
            "action":      "buy",
            "hashtag":     "Tag1",
            "quantity":    1.00,
            "user_id":     app.User.UserID,
            "bank_order":  true,
            "complete":    true,
            "value":       1.0,
            "uuid":        "complete-1",
            "created_at":  "2015-02-04T19:00:00Z",
            "executed_at": "2015-02-04T19:15:00Z",
            "resolution":  "success",
            "notes":       "",
        },
    }

    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Len(t, json_body, 2)
    json_body.Equal(t, expected)
}

func TestHistoricOrdersFilters(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    tags := []string{"Tag1", "Other"}
    for _, tagName := range tags {
        tag := models.HashTag{HashTag: tagName, Value: 1.00, InBank: 1.00}
        app.Put(tag)
    }

    orderForTag := []string{"Tag1", "Tag1", "Tag1", "Tag1", "Other", "Tag1"}
    sysOrders := []models.OrderSystem{
        models.OrderSystem{UUID: "complete-1", Complete: true, Resolution: models.SUCCESS, UserID: app.User.UserID, CreatedAt: time.Date(2015, time.February, 04, 01, 00, 0, 0, time.UTC)},
        models.OrderSystem{UUID: "complete-2", Complete: true, Resolution: models.ERROR, UserID: app.User.UserID, CreatedAt: time.Date(2015, time.February, 04, 02, 00, 0, 0, time.UTC)},
        models.OrderSystem{UUID: "complete-3", Complete: true, Resolution: models.FAILURE, UserID: app.User.UserID, CreatedAt: time.Date(2015, time.February, 04, 03, 00, 0, 0, time.UTC)},
        models.OrderSystem{UUID: "pending--1", Complete: false, Resolution: models.PENDING, UserID: app.User.UserID, CreatedAt: time.Date(2015, time.February, 04, 04, 00, 0, 0, time.UTC)},
        models.OrderSystem{UUID: "other-succ", Complete: true, Resolution: models.SUCCESS, UserID: app.User.UserID, CreatedAt: time.Date(2015, time.February, 04, 05, 00, 0, 0, time.UTC)},
        models.OrderSystem{UUID: "OTHER-USER", Complete: true, Resolution: models.SUCCESS, UserID: "stranger", CreatedAt: time.Date(2015, time.February, 04, 01, 00, 0, 0, time.UTC)},
    }

    for i, sysOrder := range sysOrders {
        order := models.Order{
            OrderBase:   models.OrderBase{HashTag: orderForTag[i]},
            OrderSystem: sysOrder,
        }
        err := app.Put(order)
        assert.NoError(t, err)
    }

    // Test Matrix
    filters := []string{"tag=Tag1", "resolution=success", "resolution=error", "tag=Tag1&resolution=success"}
    expected := [][]string{
        []string{"complete-3", "complete-2", "complete-1"},
        []string{"other-succ", "complete-1"},
        []string{"complete-2"},
        []string{"complete-1"},
    }

    for i, filter := range filters {
        // Act
        req := app.NewJsonRequest("GET", "/api/order/history/?"+filter, nil, app.User)
        rec := app.Do(req)
        json_body := app.JsonResponceToListOfStringMap(rec)

        // Assert
        assert.Equal(t, http.StatusOK, rec.Code)
        assert.Len(t, json_body, len(expected[i]))
        for j, order := range json_body {
            assert.Equal(t, order["uuid"], expected[i][j])
        }
    }
}
