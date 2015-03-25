package webapp

import (
    "net/http"
    "regexp"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
)

func apiDefinition(routes martini.Routes, r render.Render) {
    def := make(map[string]string, 0)

    exportedNameRegEx := `[A-Z]{1}.*?:[A-Z]{1}.*?`

    for _, route := range routes.All() {
        name := route.GetName()
        if matched, _ := regexp.MatchString(exportedNameRegEx, name); !matched {
            continue
        }

        def[name] = route.Pattern()
    }

    r.JSON(http.StatusOK, def)
}
