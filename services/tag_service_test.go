// Tag service
package services_test

import (
    "fmt"
    "net/http"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"

    "github.com/hashtock/hashtock-go/models"
    "github.com/hashtock/hashtock-go/test_tools"
)

func TestGetListOfTags(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/tag/", nil, app.User)

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
    app.Put(t1, t2)

    rec := app.Do(req)
    json_body := app.JsonResponceToListOfStringMap(rec)

    // Order matters
    expected := test_tools.JsonList{
        test_tools.Json{
            "hashtag": "XTag1",
            "value":   10.5,
            "in_bank": 100.0,
        },
        test_tools.Json{
            "hashtag": "Tag2",
            "value":   1,
            "in_bank": 50.2,
        },
    }

    assert.Equal(t, http.StatusOK, rec.Code)
    json_body.Equal(t, expected)
}

func TestGetSingleTag(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/tag/TestTag/", nil, app.User)

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   10.5,
        InBank:  100.0,
    }
    app.Put(tag)

    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    // Order matters
    expected := test_tools.Json{
        "hashtag": "TestTag",
        "value":   10.5,
        "in_bank": 100.0,
    }

    assert.Equal(t, http.StatusOK, rec.Code)
    json_body.Equal(t, expected)
}

func TestGetUnExistingTag(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    rec := app.ExecuteJsonRequest("GET", "/api/tag/MISSING/", nil, app.User)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "code":  http.StatusNotFound,
        "error": "HashTag \"MISSING\" not found",
    }

    assert.Equal(t, http.StatusNotFound, rec.Code)
    json_body.Equal(t, expected) // This is not very robust for error msg
}

func TestValuesForUnExistingTag(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    rec := app.ExecuteJsonRequest("GET", "/api/tag/MISSING/values/", nil, app.User)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "code":  http.StatusNotFound,
        "error": "Tag 'MISSING' does not exist",
    }

    assert.Equal(t, http.StatusNotFound, rec.Code)
    json_body.Equal(t, expected)
}

func TestValuesExistingTagNoValuesYet(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/tag/TestTag/values/", nil, app.User)
    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  100.0,
    }
    app.Put(tag)

    rec := app.Do(req)

    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "null", rec.Body.String())
}

func TestValuesForTag(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/tag/TestTag/values/", nil, app.User)

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  100.0,
    }
    app.Put(tag)

    time_1 := time.Now().Add(time.Minute * -2).Truncate(time.Second)
    tagValue := models.HashTagValue{
        HashTag: "TestTag",
        Value:   30,
        Date:    time_1,
    }
    app.Put(tagValue)

    time_2 := time.Now().Add(time.Minute * -1).Truncate(time.Second)
    tagValue2 := models.HashTagValue{
        HashTag: "TestTag",
        Value:   29,
        Date:    time_2,
    }
    app.Put(tagValue2)

    expected := test_tools.JsonList{
        test_tools.Json{
            "value": 30,
            "date":  time_1.Format(time.RFC3339),
        },
        test_tools.Json{
            "value": 29,
            "date":  time_2.Format(time.RFC3339),
        },
    }

    rec := app.Do(req)
    json_body := app.JsonResponceToListOfStringMap(rec)

    assert.Equal(t, http.StatusOK, rec.Code)
    json_body.Equal(t, expected)
}

func TestValuesForTagSamplingInvalid(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/tag/TestTag/values/?days=2", nil, app.User)

    rec := app.Do(req)
    json_body := app.JsonResponceToStringMap(rec)

    expected := test_tools.Json{
        "code":  400,
        "error": "Duration 2 not supported",
    }

    assert.Equal(t, http.StatusBadRequest, rec.Code)
    json_body.Equal(t, expected)
}

func TestValuesForTagSamplingValid(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    validDurations := []string{
        "", // 1
        "1", "7", "14", "30",
    }

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  100.0,
    }
    app.Put(tag)

    for _, duration := range validDurations {
        req := app.NewJsonRequest("GET", fmt.Sprintf("/api/tag/TestTag/values/?days=%v", duration), nil, app.User)

        rec := app.Do(req)
        assert.Equal(t, http.StatusOK, rec.Code)
    }
}

func TestValuesForTagSampling7Days(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    req := app.NewJsonRequest("GET", "/api/tag/TestTag/values/?days=7", nil, app.User)

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   1,
        InBank:  100.0,
    }
    app.Put(tag)

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
            Date:    time.Now().Truncate(24 * time.Hour).Add(offset),
        }
        err := app.Put(tagValue)
        assert.NoError(t, err)
    }

    rec := app.Do(req)
    json_body := app.JsonResponceToListOfStringMap(rec)

    expected := test_tools.JsonList{
        test_tools.Json{
            "value": 1,
            "date":  time.Now().Add(-7 * 24 * time.Hour).Truncate(24 * time.Hour).Format(time.RFC3339),
        },
        test_tools.Json{
            "value": (1 + 2) / 2.0,
            "date":  time.Now().Add(-6 * 24 * time.Hour).Truncate(24 * time.Hour).Format(time.RFC3339),
        },
        test_tools.Json{
            "value": 2,
            "date":  time.Now().Add(-5 * 24 * time.Hour).Truncate(24 * time.Hour).Format(time.RFC3339),
        },
        test_tools.Json{
            "value": (6 + 4) / 2.0,
            "date":  time.Now().Add(-1 * 24 * time.Hour).Truncate(24 * time.Hour).Format(time.RFC3339),
        },
    }

    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Len(t, json_body, 4)
    json_body.Equal(t, expected)
}
