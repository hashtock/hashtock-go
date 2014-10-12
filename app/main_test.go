package main

import (
    "encoding/json"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "appengine/aetest"

    "github.com/stretchr/testify/suite"
)

type ExampleTestSuite struct {
    suite.Suite

    inst aetest.Instance
}

type Json map[string]interface{}

/* Suite setup and generic helper methods */

func (s *ExampleTestSuite) SetupTest() {
    var err error
    s.inst, err = aetest.NewInstance(nil)
    if err != nil {
        s.T().Fatal(err)
    }
}

func (s *ExampleTestSuite) TearDownTest() {
    s.inst.Close()
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

func (s *ExampleTestSuite) ExecuteJsonRequest(method, urlStr string, body io.Reader) (rec *httptest.ResponseRecorder) {
    req := s.NewRequest("GET", "/api/", nil)
    req.Header.Add("Accept", "application/json")
    return s.Do(req)
}

func (s *ExampleTestSuite) jsonResponceToStringMap(rec *httptest.ResponseRecorder) Json {
    json_map := Json{}
    json.Unmarshal(rec.Body.Bytes(), &json_map)
    return json_map
}

/* Actuall tests */

func (s *ExampleTestSuite) TestApiHasAllEndpoints() {
    rec := s.ExecuteJsonRequest("GET", "/api/", nil)
    json_body := s.jsonResponceToStringMap(rec)

    expected := Json{
        "user": "/api/user/",
    }

    s.Equal(http.StatusOK, rec.Code)
    s.Equal(expected, json_body)
}

/* Kickoff Test Suite */

func TestExampleTestSuite(t *testing.T) {
    suite.Run(t, new(ExampleTestSuite))
}
