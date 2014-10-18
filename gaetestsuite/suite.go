package gaetestsuite

// Kickoff Test Suite
// type ExampleTestSuite struct {
//     gaetestsuite.GAETestSuite
// }
//
// func TestExampleTestSuite(t *testing.T) {
//     suite.Run(t, new(ExampleTestSuite))
// }

import (
    "encoding/json"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "appengine/aetest"
    "appengine/user"

    "github.com/stretchr/testify/suite"
)

type GAETestSuite struct {
    suite.Suite

    Inst      aetest.Instance
    User      *user.User
    AdminUser *user.User
    NoUser    *user.User
}

type Json map[string]interface{}
type JsonList []Json

func (g *GAETestSuite) SetupTest() {
    var err error

    options := &aetest.Options{
        StronglyConsistentDatastore: true,
    }

    g.Inst, err = aetest.NewInstance(options)
    if err != nil {
        g.T().Fatal(err)
    }

    g.User = g.MakeUser()
    g.AdminUser = g.MakeAdminUser()
    g.NoUser = nil
}

func (g *GAETestSuite) TearDownTest() {
    g.Inst.Close()
}

func (g *GAETestSuite) MakeAdminUser() (u *user.User) {
    u = &user.User{
        Email: "admin@admin.prv",
        Admin: true,
    }
    return
}

func (g *GAETestSuite) MakeUser() (u *user.User) {
    u = &user.User{
        Email: "user@here.prv",
        Admin: false,
    }
    return
}

func (g *GAETestSuite) NewRequest(method, urlStr string, body io.Reader) (req *http.Request) {
    var err error
    req, err = g.Inst.NewRequest(method, urlStr, body)
    if err != nil {
        g.T().Fatal(err)
    }
    return
}

func (g *GAETestSuite) NewJsonRequest(method, urlStr string, body io.Reader, u *user.User) (req *http.Request) {
    req = g.NewRequest(method, urlStr, body)

    req.Header.Add("Accept", "application/json")

    if u != nil {
        aetest.Login(u, req)
    } else {
        aetest.Logout(req)
    }

    return
}

func (g *GAETestSuite) Do(req *http.Request) (rec *httptest.ResponseRecorder) {
    rec = httptest.NewRecorder()
    http.DefaultServeMux.ServeHTTP(rec, req)
    return
}

func (g *GAETestSuite) ExecuteJsonRequest(method, urlStr string, body io.Reader, u *user.User) (rec *httptest.ResponseRecorder) {
    req := g.NewJsonRequest(method, urlStr, body, u)

    return g.Do(req)
}

func (g *GAETestSuite) JsonResponceToStringMap(rec *httptest.ResponseRecorder) Json {
    json_map := Json{}
    json.Unmarshal(rec.Body.Bytes(), &json_map)
    return json_map
}

func (g *GAETestSuite) JsonResponceToListOfStringMap(rec *httptest.ResponseRecorder) JsonList {
    json_map := JsonList{}
    json.Unmarshal(rec.Body.Bytes(), &json_map)
    return json_map
}

func Run(t *testing.T, test_suite suite.TestingSuite) {
    suite.Run(t, test_suite)
}
