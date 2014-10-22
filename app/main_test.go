package main

import (
    "net/http"
    "testing"

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
    json_body["uuid"] = "uuid"

    expected := gaetestsuite.Json{
        "action":     "buy",
        "hashtag":    "Tag1",
        "quantity":   1.00,
        "user_id":    s.User.Email,
        "bank_order": true,
        "complete":   false,
        "uuid":       "uuid",
    }

    s.Equal(http.StatusCreated, rec.Code)
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

/* Kickoff Test Suite */

func TestFunctionalTestSuite(t *testing.T) {
    gaetestsuite.Run(t, new(FunctionalTestSuite))
}
