// Tag service
// Run as part of service test suite
package services_test

import (
    "net/http"
    "time"

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

func (s *ServicesTestSuite) TestValuesForUnExistingTag() {
    rec := s.ExecuteJsonRequest("GET", "/api/tag/MISSING/values/", nil, s.User)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "code":  http.StatusNotFound,
        "error": "Tag 'MISSING' does not exist",
    }

    s.Equal(http.StatusNotFound, rec.Code)
    s.Equal(expected, json_body)
}

func (s *ServicesTestSuite) TestValuesExistingTagNoValuesYet() {
    req := s.NewJsonRequest("GET", "/api/tag/TestTag/values/", nil, s.User)
    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  100.0,
    }
    tag.Put(req)

    rec := s.Do(req)

    s.Equal(http.StatusOK, rec.Code)
    s.Equal("null", rec.Body.String())
}

func (s *ServicesTestSuite) TestValuesForTag() {
    req := s.NewJsonRequest("GET", "/api/tag/TestTag/values/", nil, s.User)

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  100.0,
    }
    tag.Put(req)

    tagValue := models.HashTagValue{
        HashTag: "TestTag",
        Value:   30,
        Date:    time.Date(2015, time.February, 04, 19, 30, 0, 0, time.UTC),
    }
    tagValue.Put(req)

    tagValue2 := models.HashTagValue{
        HashTag: "TestTag",
        Value:   29,
        Date:    time.Date(2015, time.February, 04, 19, 15, 0, 0, time.UTC),
    }
    tagValue2.Put(req)

    expected := gaetestsuite.JsonList{
        gaetestsuite.Json{
            "value": 29,
            "date":  "2015-02-04T19:15:00Z",
        },
        gaetestsuite.Json{
            "value": 30,
            "date":  "2015-02-04T19:30:00Z",
        },
    }

    rec := s.Do(req)
    json_body := s.JsonResponceToListOfStringMap(rec)

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}
