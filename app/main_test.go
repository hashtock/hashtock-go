package main_test

import (
    "net/http"
    "testing"

    _ "github.com/hashtock/hashtock-go/app"
    "github.com/hashtock/hashtock-go/gaetestsuite"
    "github.com/hashtock/hashtock-go/models"
)

type FunctionalTestSuite struct {
    gaetestsuite.GAETestSuite
}

func (s *FunctionalTestSuite) TestApiHasAllEndpoints() {
    rec := s.ExecuteJsonRequest("GET", "/api/", nil, s.User)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "user":  "/api/user/",
        "tag":   "/api/tag/",
        "order": "/api/order/",
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *FunctionalTestSuite) TestUserHasToBeLoggedIn() {
    expectedStatus := http.StatusForbidden
    expectedBody := http.StatusText(expectedStatus) + "\n"

    rec := s.ExecuteJsonRequest("GET", "/api/user/", nil, s.NoUser)

    s.Equal(expectedBody, rec.Body.String())
    s.Equal(expectedStatus, rec.Code)
}

func (s *FunctionalTestSuite) TestProfileExistForLoggedInUser() {
    rec := s.ExecuteJsonRequest("GET", "/api/user/", nil, s.User)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "id":     "user@here.prv",
        "founds": models.StartingFounds,
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *FunctionalTestSuite) TestGetListOfTags() {
    req := s.NewJsonRequest("GET", "/api/tag/", nil, s.User)

    t1 := models.HashTag{
        HashTag: "XTag1",
        Value:   10.5,
        InBank:  100.0,
    }
    t2 := models.HashTag{
        HashTag: "Tag2",
        Value:   1,
        InBank:  50.2,
    }
    t1.Put(req)
    t2.Put(req)

    rec := s.Do(req)
    json_body := s.JsonResponceToListOfStringMap(rec)

    // Order matters
    expected := gaetestsuite.JsonList{
        gaetestsuite.Json{
            "hashtag": "XTag1",
            "value":   10.5,
            "in_bank": 100.0,
        },
        gaetestsuite.Json{
            "hashtag": "Tag2",
            "value":   1,
            "in_bank": 50.2,
        },
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *FunctionalTestSuite) TestGetSingleTag() {
    req := s.NewJsonRequest("GET", "/api/tag/TestTag/", nil, s.User)

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   10.5,
        InBank:  100.0,
    }
    tag.Put(req)

    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    // Order matters
    expected := gaetestsuite.Json{
        "hashtag": "TestTag",
        "value":   10.5,
        "in_bank": 100.0,
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *FunctionalTestSuite) TestGetUnExistingTag() {
    rec := s.ExecuteJsonRequest("GET", "/api/tag/MISSING/", nil, s.User)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "code":  http.StatusNotFound,
        "error": "HashTag \"MISSING\" not found",
    }

    s.Equal(http.StatusNotFound, rec.Code)
    s.Equal(expected, json_body) // This is not very robust for error msg
}

func (s *FunctionalTestSuite) TestAdmingAddingTag() {
    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   10.5,
    }

    body := s.ToJsonBody(tag)
    req := s.NewJsonRequest("POST", "/api/tag/", body, s.AdminUser)
    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "hashtag": "TestTag",
        "value":   10.5,
        "in_bank": 100.0,
    }

    s.Equal(http.StatusCreated, rec.Code)
    s.Equal(expected, json_body)

    tag.InBank = 100
    new_tag, err := models.GetHashTag(req, "TestTag")
    s.NoError(err)
    s.Equal(tag, *new_tag)
}

func (s *FunctionalTestSuite) TestAdmingAddingTagWithInBankValue() {
    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   10.5,
        InBank:  1.0,
    }

    body := s.ToJsonBody(tag)
    req := s.NewJsonRequest("POST", "/api/tag/", body, s.AdminUser)
    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "hashtag": "TestTag",
        "value":   10.5,
        "in_bank": 100.0,
    }

    s.Equal(http.StatusCreated, rec.Code)
    s.Equal(expected, json_body)

    new_tag, err := models.GetHashTag(req, "TestTag")
    s.NoError(err)
    s.Equal(100.0, new_tag.InBank)
}

func (s *FunctionalTestSuite) TestAdminAddingExistingTag() {
    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   10.5,
        InBank:  100.0,
    }

    body := s.ToJsonBody(tag)
    req := s.NewJsonRequest("POST", "/api/tag/", body, s.AdminUser)
    if err := tag.Put(req); err != nil {
        s.T().Fatalf(err.Error())
    }

    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "code":  400,
        "error": "Tag alread exists",
    }

    s.Equal(http.StatusBadRequest, rec.Code)
    s.Equal(expected, json_body)
}

func (s *FunctionalTestSuite) TestAdmingUpdateTagValue() {
    tmp_req := s.NewJsonRequest("GET", "/", nil, nil)
    existing_tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  50.0,
    }
    existing_tag.Put(tmp_req)

    update_tag_value := models.HashTag{
        Value: 2,
    }

    body := s.ToJsonBody(update_tag_value)
    req := s.NewJsonRequest("PUT", "/api/tag/TestTag/", body, s.AdminUser)
    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "hashtag": "TestTag",
        "value":   2,
        "in_bank": 50.0,
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *FunctionalTestSuite) TestAdmingUpdateTagValueIgnoreBankValue() {
    tmp_req := s.NewJsonRequest("GET", "/", nil, nil)
    existing_tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  50.0,
    }
    existing_tag.Put(tmp_req)

    update_tag_value := models.HashTag{
        Value:  2,
        InBank: 100.0, // This will be just ignored
    }

    body := s.ToJsonBody(update_tag_value)
    req := s.NewJsonRequest("PUT", "/api/tag/TestTag/", body, s.AdminUser)
    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "hashtag": "TestTag",
        "value":   2,
        "in_bank": 50.0,
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *FunctionalTestSuite) TestAdmingUpdateTagValueInvalidHashTag() {
    tmp_req := s.NewJsonRequest("GET", "/", nil, nil)
    existing_tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  50.0,
    }
    existing_tag.Put(tmp_req)

    update_tag_value := models.HashTag{
        HashTag: "SomethingStupid",
        Value:   2,
    }

    body := s.ToJsonBody(update_tag_value)
    req := s.NewJsonRequest("PUT", "/api/tag/TestTag/", body, s.AdminUser)
    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "code":  400,
        "error": "hashtag value has to be empty or correct",
    }
    tag, _ := models.GetHashTag(req, "TestTag")

    s.Equal(http.StatusBadRequest, rec.Code)
    s.Equal(expected, json_body)
    s.Equal(existing_tag, *tag) // No change!!
}

func (s *FunctionalTestSuite) TestAdmingUpdateTagInvalidValue() {
    tmp_req := s.NewJsonRequest("GET", "/", nil, nil)
    existing_tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  50.0,
    }
    existing_tag.Put(tmp_req)

    update_tag_value := models.HashTag{
        Value: 0,
    }

    body := s.ToJsonBody(update_tag_value)
    req := s.NewJsonRequest("PUT", "/api/tag/TestTag/", body, s.AdminUser)
    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "code":  400,
        "error": "Value has to be positive",
    }
    tag, _ := models.GetHashTag(req, "TestTag")

    s.Equal(http.StatusBadRequest, rec.Code)
    s.Equal(expected, json_body)
    s.Equal(existing_tag, *tag) // No change!!
}

func (s *FunctionalTestSuite) TestRegularUserUpdateTagValue() {
    tmp_req := s.NewJsonRequest("GET", "/", nil, nil)
    existing_tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  50.0,
    }
    existing_tag.Put(tmp_req)

    update_tag_value := models.HashTag{
        Value: 2,
    }

    body := s.ToJsonBody(update_tag_value)
    req := s.NewJsonRequest("PUT", "/api/tag/TestTag/", body, s.User)
    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "code":  http.StatusForbidden,
        "error": "Forbidden",
    }
    tag, _ := models.GetHashTag(req, "TestTag")

    s.Equal(http.StatusForbidden, rec.Code)
    s.Equal(expected, json_body)
    s.Equal(existing_tag, *tag) // No change!!
}

func (s *FunctionalTestSuite) TestRegularUserAddingTag() {
    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   10.5,
    }

    body := s.ToJsonBody(tag)
    req := s.NewJsonRequest("POST", "/api/tag/", body, s.User)
    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "code":  http.StatusForbidden,
        "error": "Forbidden",
    }

    s.Equal(http.StatusForbidden, rec.Code)
    s.Equal(expected, json_body)
}

func (s *FunctionalTestSuite) TestGetUsersTags() {
    req := s.NewJsonRequest("GET", "/api/user/tags/", nil, s.User)

    t1 := models.TagShare{HashTag: "Tag1", Quantity: 10.5, UserID: s.User.Email}
    t2 := models.TagShare{HashTag: "Tag2", Quantity: 0.20, UserID: s.User.Email}
    t3 := models.TagShare{HashTag: "Tag1", Quantity: 1.00, UserID: "OtherID"}
    t1.Put(req)
    t2.Put(req)
    t3.Put(req)

    rec := s.Do(req)
    json_body := s.JsonResponceToListOfStringMap(rec)

    // Order matters
    expected := gaetestsuite.JsonList{
        gaetestsuite.Json{
            "hashtag":  "Tag1",
            "quantity": 10.5,
            "user_id":  s.User.Email,
        },
        gaetestsuite.Json{
            "hashtag":  "Tag2",
            "quantity": 0.2,
            "user_id":  s.User.Email,
        },
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *FunctionalTestSuite) TestPlaceTransactionOrderWithBank() {
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
    expected := gaetestsuite.Json{
        "action":     "buy",
        "hashtag":    "Tag1",
        "quantity":   1.00,
        "user_id":    s.User.Email,
        "bank_order": true,
        "complete":   false,
        "uuid":       uuid.(string),
    }
    order_in_db, err := models.GetOrder(req, uuid.(string))
    s.NoError(err)
    s.NotNil(order_in_db)
    s.Equal(expected, json_body)
}

func (s *FunctionalTestSuite) TestPlaceInvalidTransactionOrderWithBank() {
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

func (s *FunctionalTestSuite) TestGetOrderDetails() {
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
            UserID:   s.User.Email,
            Complete: false,
        },
    }
    order.Put(req)

    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "action":     "buy",
        "hashtag":    "Tag1",
        "quantity":   1.00,
        "user_id":    s.User.Email,
        "bank_order": true,
        "complete":   false,
        "uuid":       "FAKE-UUID",
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *FunctionalTestSuite) TestGetOrderDifferentUser() {
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

func (s *FunctionalTestSuite) TestCancelOrder() {
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
    s.Equal(0, rec.Body.Len())

    cancelled_order, _ := models.GetOrder(req, "FAKE-UUID")
    s.Nil(cancelled_order)
}

func (s *FunctionalTestSuite) TestCancelCompetedOrder() {
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

func (s *FunctionalTestSuite) TestCurrentOrders() {
    req := s.NewJsonRequest("GET", "/api/order/", nil, s.User)

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
            UUID:     "pending-1",
            UserID:   s.User.Email,
            Complete: false,
        },
    }
    order_1.Put(req)

    order_1.UUID = "pending-2"
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
            "action":     "buy",
            "hashtag":    "Tag1",
            "quantity":   1.00,
            "user_id":    s.User.Email,
            "bank_order": true,
            "complete":   false,
            "uuid":       "pending-1",
        },
        gaetestsuite.Json{
            "action":     "buy",
            "hashtag":    "Tag1",
            "quantity":   1.00,
            "user_id":    s.User.Email,
            "bank_order": true,
            "complete":   false,
            "uuid":       "pending-2",
        },
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Len(json_body, 2)
    s.Equal(expected, json_body)
}

func (s *FunctionalTestSuite) TestHistoricOrders() {
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
            UUID:     "complete-1",
            UserID:   s.User.Email,
            Complete: true,
        },
    }
    order_1.Put(req)

    order_1.UUID = "complete-2"
    order_1.Put(req)

    order_2 := models.Order{
        OrderBase: base_order,
        OrderSystem: models.OrderSystem{
            UUID:     "pending-1",
            UserID:   s.User.Email,
            Complete: false,
        },
    }
    order_2.Put(req)

    order_3 := models.Order{
        OrderBase: base_order,
        OrderSystem: models.OrderSystem{
            UUID:     "complete-3",
            UserID:   "some user",
            Complete: true,
        },
    }
    order_3.Put(req)

    rec := s.Do(req)
    json_body := s.JsonResponceToListOfStringMap(rec)

    expected := gaetestsuite.JsonList{
        gaetestsuite.Json{
            "action":     "buy",
            "hashtag":    "Tag1",
            "quantity":   1.00,
            "user_id":    s.User.Email,
            "bank_order": true,
            "complete":   true,
            "uuid":       "complete-1",
        },
        gaetestsuite.Json{
            "action":     "buy",
            "hashtag":    "Tag1",
            "quantity":   1.00,
            "user_id":    s.User.Email,
            "bank_order": true,
            "complete":   true,
            "uuid":       "complete-2",
        },
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Len(json_body, 2)
    s.Equal(expected, json_body)
}

/* Kickoff Test Suite */

func TestFunctionalTestSuite(t *testing.T) {
    gaetestsuite.Run(t, new(FunctionalTestSuite))
}
