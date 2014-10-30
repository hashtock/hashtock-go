// Current user service
// Run as part of service test suite
package services_test

import (
    "net/http"

    "github.com/hashtock/hashtock-go/gaetestsuite"
    "github.com/hashtock/hashtock-go/models"
)

// TODO(security): It would be good to expand this to test ALL api urls
func (s *ServicesTestSuite) TestUserHasToBeLoggedIn() {
    expectedStatus := http.StatusForbidden
    expectedBody := http.StatusText(expectedStatus) + "\n"

    rec := s.ExecuteJsonRequest("GET", "/api/user/", nil, s.NoUser)

    s.Equal(expectedBody, rec.Body.String())
    s.Equal(expectedStatus, rec.Code)
}

func (s *ServicesTestSuite) TestProfileExistForLoggedInUser() {
    rec := s.ExecuteJsonRequest("GET", "/api/user/", nil, s.User)
    json_body := s.JsonResponceToStringMap(rec)

    expected := gaetestsuite.Json{
        "id":     "user@here.prv",
        "founds": models.StartingFounds,
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *ServicesTestSuite) TestGetUsersTags() {
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
