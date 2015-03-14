// Tag service
// Run as part of service test suite
package services_test

import (
    "fmt"
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

    time_1 := time.Now().Add(time.Minute * -2).Truncate(time.Second)
    tagValue := models.HashTagValue{
        HashTag: "TestTag",
        Value:   30,
        Date:    time_1,
    }
    tagValue.Put(req)

    time_2 := time.Now().Add(time.Minute * -1).Truncate(time.Second)
    tagValue2 := models.HashTagValue{
        HashTag: "TestTag",
        Value:   29,
        Date:    time_2,
    }
    tagValue2.Put(req)

    expected := gaetestsuite.JsonList{
        gaetestsuite.Json{
            "value": 30,
            "date":  time_1.Format(time.RFC3339),
        },
        gaetestsuite.Json{
            "value": 29,
            "date":  time_2.Format(time.RFC3339),
        },
    }

    rec := s.Do(req)
    json_body := s.JsonResponceToListOfStringMap(rec)

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *ServicesTestSuite) TestValuesForTagSamplingInvalid() {
    req := s.NewJsonRequest("GET", "/api/tag/TestTag/values/?days=2", nil, s.User)

    rec := s.Do(req)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "code":  400,
        "error": "Duration 2 not supported",
    }

    s.Equal(http.StatusBadRequest, rec.Code)
    s.Equal(expected, json_body)
}

func (s *ServicesTestSuite) TestValuesForTagSamplingValid() {
    validDurations := []string{
        "", // 1
        "1", "7", "14", "30",
    }

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  100.0,
    }
    tag.Put(s.DummyRequest(s.User))

    for _, duration := range validDurations {
        req := s.NewJsonRequest("GET", fmt.Sprintf("/api/tag/TestTag/values/?days=%v", duration), nil, s.User)

        rec := s.Do(req)
        s.Equal(http.StatusOK, rec.Code)
    }
}

func (s *ServicesTestSuite) TestValuesForTagSampling7Days() {
    req := s.NewJsonRequest("GET", "/api/tag/TestTag/values/?days=7", nil, s.User)

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  100.0,
    }
    tag.Put(req)

    tagValuesOffsets := map[time.Duration]float64{
        -7*24*time.Hour + 5*time.Minute: 1,
        -6*24*time.Hour + 1*time.Hour:   1,
        -6*24*time.Hour + 2*time.Hour:   2,
        -5*24*time.Hour + 1*time.Hour:   2,
        -2*time.Hour + 15*time.Minute:   4,
        -1*time.Hour + 30*time.Minute:   6,
    }

    for offset, value := range tagValuesOffsets {
        tagValue := models.HashTagValue{
            HashTag: "TestTag",
            Value:   value,
            Date:    time.Now().Add(offset),
        }
        tagValue.Put(req)
    }

    rec := s.Do(req)
    json_body := s.JsonResponceToListOfStringMap(rec)

    expected := gaetestsuite.JsonList{
        gaetestsuite.Json{
            "value": 1,
            "date":  time.Now().Add(-7 * 24 * time.Hour).Truncate(24 * time.Hour).Format(time.RFC3339),
        },
        gaetestsuite.Json{
            "value": (1 + 2) / 2.0,
            "date":  time.Now().Add(-6 * 24 * time.Hour).Truncate(24 * time.Hour).Format(time.RFC3339),
        },
        gaetestsuite.Json{
            "value": 2,
            "date":  time.Now().Add(-5 * 24 * time.Hour).Truncate(24 * time.Hour).Format(time.RFC3339),
        },
        gaetestsuite.Json{
            "value": (6 + 4) / 2.0,
            "date":  time.Now().Truncate(24 * time.Hour).Format(time.RFC3339),
        },
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Len(json_body, 4)
    s.JsonListEqual(expected, json_body)
}
