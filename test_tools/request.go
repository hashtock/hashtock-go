package test_tools

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "net/http/httptest"

    "github.com/gorilla/context"

    "github.com/hashtock/hashtock-go/models"
)

func (t *TestApp) Login(req *http.Request, profile *models.Profile) {
    context.Set(req, "reqProfile", profile)
}

func (t *TestApp) ToJsonBody(data interface{}) (body io.Reader) {
    marshaled, err := json.Marshal(data)
    if err != nil {
        t.test.Fatal(err)
    }
    return bytes.NewBuffer(marshaled)
}

func (t *TestApp) NewJsonRequest(method, urlStr string, body io.Reader, p *models.Profile) (req *http.Request) {
    var err error
    req, err = http.NewRequest(method, urlStr, body)
    if err != nil {
        t.test.Fatal(err)
    }
    t.Login(req, p)

    req.Header.Add("Accept", "application/json")
    if method == "POST" || method == "PUT" {
        req.Header.Add("content-type", "application/json")
    }

    return
}

func (t *TestApp) JsonResponceToStringMap(rec *httptest.ResponseRecorder) Json {
    json_map := Json{}

    if rec.Body.Len() == 0 {
        t.test.Fatalf("Response is empty!")
    }

    if err := json.Unmarshal(rec.Body.Bytes(), &json_map); err != nil {
        t.test.Fatalf("Could not unmarshal map: %s. Err: %v. Headers: %#v", rec.Body.String(), err, rec.HeaderMap)
    }
    return json_map
}

func (t *TestApp) JsonResponceToListOfStringMap(rec *httptest.ResponseRecorder) JsonList {
    json_map := JsonList{}
    if rec.Body.Len() == 0 {
        t.test.Fatalf("Response is empty!")
    }

    if err := json.Unmarshal(rec.Body.Bytes(), &json_map); err != nil {
        t.test.Fatalf("Could not unmarshal list: %s", rec.Body.String())
    }
    return json_map
}

func (t *TestApp) ExecuteJsonRequestRaw(method, urlStr string, data interface{}, p *models.Profile) (rec *httptest.ResponseRecorder) {
    body := t.ToJsonBody(data)
    return t.ExecuteJsonRequest(method, urlStr, body, p)
}

func (t *TestApp) ExecuteJsonRequest(method, urlStr string, body io.Reader, p *models.Profile) (rec *httptest.ResponseRecorder) {
    req := t.NewJsonRequest(method, urlStr, body, p)

    return t.Do(req)
}

func (t *TestApp) Do(req *http.Request) (rec *httptest.ResponseRecorder) {
    rec = httptest.NewRecorder()
    t.server.Config.Handler.ServeHTTP(rec, req)
    return
}
