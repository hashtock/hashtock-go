package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/gorilla/mux"

    "github.com/hashtock/hashtock-go/http_utils"
)

type Resourcer interface {
    Name() string
    EndPoints() []*EndPoint
}

type Api struct {
    endpoints map[string]string
    resources map[string]Resourcer
    router    *mux.Router
    baseURI   string
}

func NewApi(router *mux.Router, resource ...Resourcer) *Api {
    api := Api{ //a.resource_path(resource, endpoint, true)
        endpoints: map[string]string{},
        // endpoints: map[string]map[string]*EndPoint{},
    }

    api.SetRouter(router)

    for _, res := range resource {
        if err := api.AddResouce(res); err != nil {
            log.Println(err)
        }
    }

    return &api
}

func (a *Api) SetRouter(router *mux.Router) {
    a.router = router
    a.router.HandleFunc("/", a.resouce_endpoints).Methods("GET").Name("API_INFO")

    base_route := a.router.Get("API_INFO")
    if url, err := base_route.URLPath(); err == nil {
        a.baseURI = url.Path
    } else {
        log.Printf("Could not get base URI for Api. %s", err)
    }
}

func (a *Api) AddResouce(resource Resourcer) error {
    name := resource.Name()
    if _, exist := a.resources[name]; exist {
        return fmt.Errorf("API: Resource %#v already on added", resource)
    }

    // a.endpoints[name] = map[string]*EndPoint{}

    a.registerEndpoints(resource)

    return nil
}

func (a *Api) MarshalJSON() ([]byte, error) {
    return json.Marshal(a.endpoints)
}

func (a *Api) resource_path(resource Resourcer, endpoint *EndPoint, full bool) (path string) {
    var (
        buff                  bytes.Buffer
        baseURI, endpoint_uri string
    )

    if endpoint != nil {
        endpoint_uri = endpoint.URI
    }

    if full {
        baseURI = a.baseURI
    }

    url_parts := []string{
        strings.Trim(baseURI, "/"),
        strings.Trim(resource.Name(), "/"),
        strings.Trim(endpoint_uri, "/"),
    }

    buff.WriteString("/")
    for _, str := range url_parts {
        if str == "" {
            continue
        }

        buff.WriteString(strings.Trim(str, "/"))
        buff.WriteString("/")
    }
    path = buff.String()

    return
}

func (a *Api) registerEndpoints(resource Resourcer) {
    uri_prefix := a.resource_path(resource, nil, false)
    subrouter := a.router.PathPrefix(uri_prefix).Subrouter()

    for _, endpoint := range resource.EndPoints() {
        subrouter.HandleFunc(endpoint.URI, endpoint.Handler).Methods(endpoint.Method).Name(endpoint.Name)

        if endpoint.isMain() {
            a.endpoints[resource.Name()] = a.resource_path(resource, endpoint, true)
        }
    }
}

func (a *Api) resouce_endpoints(rw http.ResponseWriter, req *http.Request) {
    http_utils.SerializeResponse(rw, req, a, http.StatusOK)
}
