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

func (s *ExampleTestSuite) Do(req *http.Request) (rec *httptest.ResponseRecorder) {
    rec = httptest.NewRecorder()
    http.DefaultServeMux.ServeHTTP(rec, req)
    return
}

func (s *ExampleTestSuite) ExecuteJsonRequest(method, urlStr string, body io.Reader, user *user.User) (rec *httptest.ResponseRecorder) {
    req := s.NewRequest("GET", urlStr, nil)
    req.Header.Add("Accept", "application/json")

    if user != nil {
        aetest.Login(user, req)
    } else {
        aetest.Logout(req)
    }

    return s.Do(req)
}

func (s *ExampleTestSuite) jsonResponceToStringMap(rec *httptest.ResponseRecorder) Json {
    json_map := Json{}
    json.Unmarshal(rec.Body.Bytes(), &json_map)
    return json_map
}

/* Actuall tests */

func (s *ExampleTestSuite) TestApiHasAllEndpoints() {
    rec := s.ExecuteJsonRequest("GET", "/api/", nil, s.user)
    json_body := s.jsonResponceToStringMap(rec)

    expected := Json{
        "user": "/api/user/",
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

/* Kickoff Test Suite */

func TestExampleTestSuite(t *testing.T) {
    suite.Run(t, new(ExampleTestSuite))
}
