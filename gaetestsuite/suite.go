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
    "bytes"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "appengine"
    "appengine/aetest"
    "appengine/datastore"
    "appengine/user"

    "github.com/stretchr/testify/suite"
)

type GAETestSuite struct {
    suite.Suite

    Inst      aetest.Instance
    User      *user.User
    AdminUser *user.User
    NoUser    *user.User

    tm_suite_start time.Time
}

type testStruct struct {
    data string
}

type Json map[string]interface{}
type JsonList []Json

func (g *GAETestSuite) SetupSuite() {
    var err error

    options := &aetest.Options{
        StronglyConsistentDatastore: true,
    }

    g.Inst, err = aetest.NewInstance(options)
    if err != nil {
        g.T().Fatal(err)
    }

    g.tm_suite_start = time.Now()
}

func (g *GAETestSuite) TearDownSuite() {
    durration := time.Since(g.tm_suite_start)

    g.Inst.Close()

    log.Printf("Tests took: %v", durration)
}

func (g *GAETestSuite) SetupTest() {
    g.ClearDB()

    g.User = g.MakeUser()
    g.AdminUser = g.MakeAdminUser()
    g.NoUser = nil
}

func (g *GAETestSuite) ClearDB() {
    ctx := g.NewContext()
    keys, err := datastore.NewQuery("").KeysOnly().GetAll(ctx, nil)
    if err != nil {
        g.T().Fatal(err)
    }

    if err := datastore.DeleteMulti(ctx, keys); err != nil {
        g.T().Fatal(err)
    }
}

func (g *GAETestSuite) DummyRequest(u *user.User) *http.Request {
    return g.NewUserRequest("GET", "/", nil, u)
}

func (g *GAETestSuite) NewContext() appengine.Context {
    req := g.DummyRequest(nil)
    return appengine.NewContext(req)
}

func (g *GAETestSuite) datastoreSize() int {
    ctx := g.NewContext()
    count, err := datastore.NewQuery("").Count(ctx)
    if err != nil {
        g.T().Fatal(err)
    }

    return count
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

func (g *GAETestSuite) NewUserRequest(method, urlStr string, body io.Reader, u *user.User) (req *http.Request) {
    req = g.NewRequest(method, urlStr, body)

    if u != nil {
        aetest.Login(u, req)
    } else {
        aetest.Logout(req)
    }

    return
}

func (g *GAETestSuite) NewJsonRequest(method, urlStr string, body io.Reader, u *user.User) (req *http.Request) {
    req = g.NewUserRequest(method, urlStr, body, u)

    req.Header.Add("Accept", "application/json")
    if method == "POST" || method == "PUT" {
        req.Header.Add("content-type", "application/json")
    }

    return
}

func (g *GAETestSuite) ToJsonBody(data interface{}) (body io.Reader) {
    var (
        err       error
        marshaled []byte
    )

    if marshaled, err = json.Marshal(data); err != nil {
        g.T().Fatal(err)
    }
    return bytes.NewBuffer(marshaled)
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
    if rec.Body.Len() == 0 {
        g.T().Fatalf("Response is empty!")
    }

    if err := json.Unmarshal(rec.Body.Bytes(), &json_map); err != nil {
        g.T().Fatalf("Could not unmarshal map: %s", rec.Body.String())
    }
    return json_map
}

func (g *GAETestSuite) JsonResponceToListOfStringMap(rec *httptest.ResponseRecorder) JsonList {
    json_map := JsonList{}
    if rec.Body.Len() == 0 {
        g.T().Fatalf("Response is empty!")
    }

    if err := json.Unmarshal(rec.Body.Bytes(), &json_map); err != nil {
        g.T().Fatalf("Could not unmarshal list: %s", rec.Body.String())
    }
    return json_map
}

func (g *GAETestSuite) TestCleanUp() {
    ctx := g.NewContext()
    key := datastore.NewKey(ctx, "__TestEntity__", "abc", 0, nil)
    data := testStruct{}
    datastore.Put(ctx, key, &data)

    g.Equal(1, g.datastoreSize())

    g.ClearDB()

    g.Equal(0, g.datastoreSize())
}

func Run(t *testing.T, test_suite suite.TestingSuite) {
    suite.Run(t, test_suite)
}
