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
        "user": "/api/user/",
        "tag":  "/api/tag/",
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

    t1 := models.TagShare{
        HashTag:  "Tag1",
        Quantity: 10.5,
        UserID:   s.User.Email,
    }
    t2 := models.TagShare{
        HashTag:  "Tag2",
        Quantity: 0.2,
        UserID:   s.User.Email,
    }
    t3 := models.TagShare{
        HashTag:  "Tag1",
        Quantity: 1,
        UserID:   "OtherID",
    }
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

/* Kickoff Test Suite */

func TestFunctionalTestSuite(t *testing.T) {
    gaetestsuite.Run(t, new(FunctionalTestSuite))
}
