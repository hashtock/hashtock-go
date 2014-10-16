package main

import (
    "encoding/json"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "appengine/aetest"
    "appengine/user"

    "github.com/stretchr/testify/suite"

    "github.com/hashtock/hashtock-go/models"
)

type ExampleTestSuite struct {
    suite.Suite

    inst       aetest.Instance
    user       *user.User
    admin_user *user.User
    no_user    *user.User
}

type Json map[string]interface{}
type JsonList []Json

/* Suite setup and generic helper methods */

func (s *ExampleTestSuite) SetupTest() {
    var err error

    options := &aetest.Options{
        StronglyConsistentDatastore: true,
    }

    s.inst, err = aetest.NewInstance(options)
    if err != nil {
        s.T().Fatal(err)
    }

    s.user = s.User()
    s.admin_user = s.AdminUser()
    s.no_user = nil
}

func (s *ExampleTestSuite) TearDownTest() {
    s.inst.Close()
}

func (s *ExampleTestSuite) AdminUser() (u *user.User) {
    u = &user.User{
        Email: "admin@admin.prv",
        Admin: true,
    }
    return
}

func (s *ExampleTestSuite) User() (u *user.User) {
    u = &user.User{
        Email: "user@here.prv",
        Admin: false,
    }
    return
}

func (s *ExampleTestSuite) NewRequest(method, urlStr string, body io.Reader) (req *http.Request) {
    var err error
    req, err = s.inst.NewRequest(method, urlStr, body)
    if err != nil {
        s.T().Fatal(err)
    }
    return
}

func (e *ExampleTestSuite) NewJsonRequest(method, urlStr string, body io.Reader, u *user.User) (req *http.Request) {
    req = e.NewRequest(method, urlStr, body)

    req.Header.Add("Accept", "application/json")

    if u != nil {
        aetest.Login(u, req)
    } else {
        aetest.Logout(req)
    }

    return
}

func (s *ExampleTestSuite) Do(req *http.Request) (rec *httptest.ResponseRecorder) {
    rec = httptest.NewRecorder()
    http.DefaultServeMux.ServeHTTP(rec, req)
    return
}

func (s *ExampleTestSuite) ExecuteJsonRequest(method, urlStr string, body io.Reader, u *user.User) (rec *httptest.ResponseRecorder) {
    req := s.NewJsonRequest(method, urlStr, body, u)

    return s.Do(req)
}

func (s *ExampleTestSuite) jsonResponceToStringMap(rec *httptest.ResponseRecorder) Json {
    json_map := Json{}
    json.Unmarshal(rec.Body.Bytes(), &json_map)
    return json_map
}

func (s *ExampleTestSuite) jsonResponceToListOfStringMap(rec *httptest.ResponseRecorder) JsonList {
    json_map := JsonList{}
    json.Unmarshal(rec.Body.Bytes(), &json_map)
    return json_map
}

/* Actuall tests */

func (s *ExampleTestSuite) TestApiHasAllEndpoints() {
    rec := s.ExecuteJsonRequest("GET", "/api/", nil, s.user)
    json_body := s.jsonResponceToStringMap(rec)

    expected := Json{
        "user": "/api/user/",
        "tag":  "/api/tag/",
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *ExampleTestSuite) TestUserHasToBeLoggedIn() {
    expectedStatus := http.StatusForbidden
    expectedBody := http.StatusText(expectedStatus) + "\n"

    rec := s.ExecuteJsonRequest("GET", "/api/user/", nil, s.no_user)

    s.Equal(expectedBody, rec.Body.String())
    s.Equal(expectedStatus, rec.Code)
}

func (s *ExampleTestSuite) TestProfileExistForLoggedInUser() {
    rec := s.ExecuteJsonRequest("GET", "/api/user/", nil, s.user)
    json_body := s.jsonResponceToStringMap(rec)

    expected := Json{
        "id":     "user@here.prv",
        "founds": models.StartingFounds,
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *ExampleTestSuite) TestGetListOfTags() {
    req := s.NewJsonRequest("GET", "/api/tag/", nil, s.user)

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
    json_body := s.jsonResponceToListOfStringMap(rec)

    // Order matters
    expected := JsonList{
        Json{
            "hashtag": "XTag1",
            "value":   10.5,
            "in_bank": 100.0,
        },
        Json{
            "hashtag": "Tag2",
            "value":   1,
            "in_bank": 50.2,
        },
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *ExampleTestSuite) TestGetSingleTag() {
    req := s.NewJsonRequest("GET", "/api/tag/TestTag/", nil, s.user)

    tag := models.HashTag{
        HashTag: "TestTag",
        Value:   10.5,
        InBank:  100.0,
    }
    tag.Put(req)

    rec := s.Do(req)
    json_body := s.jsonResponceToStringMap(rec)

    // Order matters
    expected := Json{
        "hashtag": "TestTag",
        "value":   10.5,
        "in_bank": 100.0,
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

func (s *ExampleTestSuite) TestGetUnExistingTag() {
    rec := s.ExecuteJsonRequest("GET", "/api/tag/MISSING/", nil, s.user)
    json_body := s.jsonResponceToStringMap(rec)

    expected := Json{
        "code":  http.StatusNotFound,
        "error": "HashTag \"MISSING\" not found",
    }

    s.Equal(http.StatusNotFound, rec.Code)
    s.Equal(expected, json_body) // This is not very robust for error msg
}

/* Kickoff Test Suite */

func TestExampleTestSuite(t *testing.T) {
    suite.Run(t, new(ExampleTestSuite))
}
