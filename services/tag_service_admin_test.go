// Admin part of Tag service (candidate for separate service?)
// Run as part of service test suite
package services_test

import (
    "net/http"

    "github.com/hashtock/hashtock-go/gaetestsuite"
    "github.com/hashtock/hashtock-go/models"
)

func (s *ServicesTestSuite) TestAdminAddingTag() {
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

func (s *ServicesTestSuite) TestAdminAddingTagWithInBankValue() {
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

func (s *ServicesTestSuite) TestAdminAddingExistingTag() {
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

func (s *ServicesTestSuite) TestAdminUpdateTagValue() {
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

func (s *ServicesTestSuite) TestAdminUpdateTagValueIgnoreBankValue() {
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

func (s *ServicesTestSuite) TestAdminUpdateTagValueInvalidHashTag() {
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

func (s *ServicesTestSuite) TestAdminUpdateTagInvalidValue() {
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

func (s *ServicesTestSuite) TestRegularUserUpdateTagValue() {
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

func (s *ServicesTestSuite) TestRegularUserAddingTag() {
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
