// Tag service
// Run as part of service test suite
package services_test

import (
    "net/http"

    "github.com/hashtock/hashtock-go/gaetestsuite"
    "github.com/hashtock/hashtock-go/models"
)

func (s *ServicesTestSuite) TestGetListOfTags() {
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

func (s *ServicesTestSuite) TestGetSingleTag() {
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

func (s *ServicesTestSuite) TestGetUnExistingTag() {
    rec := s.ExecuteJsonRequest("GET", "/api/tag/MISSING/", nil, s.User)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "code":  http.StatusNotFound,
        "error": "HashTag \"MISSING\" not found",
    }

    s.Equal(http.StatusNotFound, rec.Code)
    s.Equal(expected, json_body) // This is not very robust for error msg
}
