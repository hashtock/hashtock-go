// Order service
// Run as part of service test suite
package services_test

import (
    "net/http"
    "time"

    "github.com/hashtock/hashtock-go/gaetestsuite"
    "github.com/hashtock/hashtock-go/models"
)

func (s *ServicesTestSuite) TestPlaceTransactionOrderWithBank() {
    order := models.OrderBase{
        Action:    "buy",
        BankOrder: true,
        HashTag:   "Tag1",
        Quantity:  1.00,
    }

    body := s.ToJsonBody(order)
    req := s.NewJsonRequest("POST", "/api/order/", body, s.User)

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    tag.Put(req)

    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)
    s.Equal(http.StatusCreated, rec.Code)
    if http.StatusCreated != rec.Code {
        // There is no point of going further!
        s.T().Fatalf("Code incorrect. Body: %v", rec.Body.String())
    }

    uuid := json_body["uuid"]
    created := json_body["created_at"]
    expected := gaetestsuite.Json{
        "action":      "buy",
        "hashtag":     "Tag1",
        "quantity":    1.00,
        "user_id":     s.User.Email,
        "bank_order":  true,
        "complete":    false,
        "value":       1.00,
        "uuid":        uuid.(string),
        "created_at":  created,
        "executed_at": "0001-01-01T00:00:00Z", // Not executed yet
        "resolution":  "",
        "notes":       "",
    }
    order_in_db, err := models.GetOrder(req, uuid.(string))
    s.NoError(err)
    s.NotNil(order_in_db)
    s.JsonEqual(expected, json_body)
}

func (s *ServicesTestSuite) TestPlaceInvalidTransactionOrderWithBank() {
    order := models.OrderBase{
        Action:    "freebe",
        BankOrder: true,
        HashTag:   "",
        Quantity:  101.00,
    }

    body := s.ToJsonBody(order)
    req := s.NewJsonRequest("POST", "/api/order/", body, s.User)

    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "code":  400,
        "error": "Incorrect fields: action, hashtag, quantity",
    }

    s.Equal(http.StatusBadRequest, rec.Code)
    s.Equal(expected, json_body)
}

func (s *ServicesTestSuite) TestGetOrderDetails() {
    req := s.NewJsonRequest("GET", "/api/order/FAKE-UUID/", nil, s.User)

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    tag.Put(req)

    order := models.Order{
        OrderBase: models.OrderBase{
            Action:    "buy",
            BankOrder: true,
            HashTag:   "Tag1",
            Quantity:  1.00,
        },
        OrderSystem: models.OrderSystem{
            UUID:       "FAKE-UUID",
            UserID:     s.User.Email,
            Complete:   false,
            Value:      1,
            CreatedAt:  time.Date(2015, time.February, 04, 19, 30, 0, 0, time.UTC),
            ExecutedAt: time.Time{},
            Resolution: models.PENDING,
            Notes:      "Some note",
        },
    }
    time.Local = time.UTC
    order.Put(req)

    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "action":      "buy",
        "hashtag":     "Tag1",
        "quantity":    1.00,
        "user_id":     s.User.Email,
        "bank_order":  true,
        "complete":    false,
        "value":       1.00,
        "uuid":        "FAKE-UUID",
        "created_at":  "2015-02-04T19:30:00Z",
        "executed_at": "0001-01-01T00:00:00Z", // Not executed yet
        "resolution":  "",
        "notes":       "Some note",
    }

    s.Equal(http.StatusOK, rec.Code)
    s.JsonEqual(expected, json_body)
}

func (s *ServicesTestSuite) TestGetOrderDifferentUser() {
    req := s.NewJsonRequest("GET", "/api/order/FAKE-UUID/", nil, s.User)

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    tag.Put(req)

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
    order.Put(req)

    rec := s.Do(req)

    s.Equal(http.StatusNotFound, rec.Code)
}

func (s *ServicesTestSuite) TestCancelOrder() {
    req := s.NewJsonRequest("DELETE", "/api/order/FAKE-UUID/", nil, s.User)

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    tag.Put(req)

    order := models.Order{
        OrderBase: models.OrderBase{
            Action:    "buy",
            BankOrder: true,
            HashTag:   "Tag1",
            Quantity:  1.00,
        },
        OrderSystem: models.OrderSystem{
            UUID:     "FAKE-UUID",
            UserID:   s.User.Email,
            Complete: false,
        },
    }
    order.Put(req)

    rec := s.Do(req)

    s.Equal(http.StatusNoContent, rec.Code)
    s.Equal(0, rec.Body.Len(), rec.Body.String())

    cancelled_order, _ := models.GetOrder(req, "FAKE-UUID")
    s.Nil(cancelled_order)
}

func (s *ServicesTestSuite) TestCancelCompetedOrder() {
    req := s.NewJsonRequest("DELETE", "/api/order/FAKE-UUID/", nil, s.User)

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    tag.Put(req)

    order := models.Order{
        OrderBase: models.OrderBase{
            Action:    "buy",
            BankOrder: true,
            HashTag:   "Tag1",
            Quantity:  1.00,
        },
        OrderSystem: models.OrderSystem{
            UUID:     "FAKE-UUID",
            UserID:   s.User.Email,
            Complete: true,
        },
    }
    order.Put(req)

    rec := s.Do(req)

    s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *ServicesTestSuite) TestCurrentOrders() {
    req := s.NewJsonRequest("GET", "/api/order/", nil, s.User)
    time.Local = time.UTC

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    tag.Put(req)

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
            UserID:     s.User.Email,
            Complete:   false,
            Value:      1,
            CreatedAt:  time.Date(2015, time.February, 04, 19, 30, 0, 0, time.UTC),
            ExecutedAt: time.Time{},
        },
    }
    order_1.Put(req)

    order_1.UUID = "pending-2"
    order_1.CreatedAt = time.Date(2015, time.February, 04, 19, 31, 0, 0, time.UTC)
    order_1.ExecutedAt = time.Time{}
    order_1.Put(req)

    order_2 := models.Order{
        OrderBase: base_order,
        OrderSystem: models.OrderSystem{
            UUID:     "complete-1",
            UserID:   s.User.Email,
            Complete: true,
        },
    }
    order_2.Put(req)

    order_3 := models.Order{
        OrderBase: base_order,
        OrderSystem: models.OrderSystem{
            UUID:     "pending-3",
            UserID:   "some user",
            Complete: false,
        },
    }
    order_3.Put(req)

    rec := s.Do(req)
    json_body := s.JsonResponceToListOfStringMap(rec)

    expected := gaetestsuite.JsonList{
        gaetestsuite.Json{
            "action":      "buy",
            "hashtag":     "Tag1",
            "quantity":    1.00,
            "user_id":     s.User.Email,
            "bank_order":  true,
            "complete":    false,
            "uuid":        "pending-2",
            "value":       1.00,
            "created_at":  "2015-02-04T19:31:00Z",
            "executed_at": "0001-01-01T00:00:00Z", // Not executed yet
            "resolution":  "",
            "notes":       "",
        },
        gaetestsuite.Json{
            "action":      "buy",
            "hashtag":     "Tag1",
            "quantity":    1.00,
            "user_id":     s.User.Email,
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

    s.Equal(http.StatusOK, rec.Code)
    s.Len(json_body, 2)
    s.JsonListEqual(expected, json_body)
}

func (s *ServicesTestSuite) TestHistoricOrders() {
    req := s.NewJsonRequest("GET", "/api/order/history/", nil, s.User)

    tag := models.HashTag{HashTag: "Tag1", Value: 1.00, InBank: 1.00}
    tag.Put(req)

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
            UserID:     s.User.Email,
            Complete:   true,
            Value:      1.0,
            CreatedAt:  time.Date(2015, time.February, 04, 19, 00, 0, 0, time.UTC),
            ExecutedAt: time.Date(2015, time.February, 04, 19, 15, 0, 0, time.UTC),
            Resolution: models.SUCCESS,
        },
    }
    order_1.Put(req)

    order_1.UUID = "complete-2"
    order_1.CreatedAt = time.Date(2015, time.February, 04, 19, 30, 0, 0, time.UTC)
    order_1.ExecutedAt = time.Date(2015, time.February, 04, 19, 45, 0, 0, time.UTC)
    order_1.Resolution = models.FAILURE
    order_1.Notes = "some reason"
    order_1.Put(req)

    order_2 := models.Order{
        OrderBase: base_order,
        OrderSystem: models.OrderSystem{
            UUID:       "pending-1",
            UserID:     s.User.Email,
            Complete:   false,
            Value:      1.0,
            CreatedAt:  time.Date(2015, time.February, 04, 19, 00, 0, 0, time.UTC),
            ExecutedAt: time.Time{},
        },
    }
    order_2.Put(req)

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
    order_3.Put(req)

    rec := s.Do(req)
    json_body := s.JsonResponceToListOfStringMap(rec)

    expected := gaetestsuite.JsonList{
        gaetestsuite.Json{
            "action":      "buy",
            "hashtag":     "Tag1",
            "quantity":    1.00,
            "user_id":     s.User.Email,
            "bank_order":  true,
            "complete":    true,
            "value":       1.0,
            "uuid":        "complete-2",
            "created_at":  "2015-02-04T19:30:00Z",
            "executed_at": "2015-02-04T19:45:00Z",
            "resolution":  "failure",
            "notes":       "some reason",
        },
        gaetestsuite.Json{
            "action":      "buy",
            "hashtag":     "Tag1",
            "quantity":    1.00,
            "user_id":     s.User.Email,
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

    s.Equal(http.StatusOK, rec.Code)
    s.Len(json_body, 2)
    s.JsonListEqual(expected, json_body)
}
