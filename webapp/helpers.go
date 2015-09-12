package webapp

import (
	"errors"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/schema"
	authCore "github.com/hashtock/auth/core"

	"github.com/hashtock/hashtock-go/core"
)

func userId(req *http.Request) (id string, err error) {
	obj := context.Get(req, UserContextKey)
	if obj == nil {
		err = errors.New("User not logged in")
		return
	}

	user, ok := obj.(*authCore.User)
	if !ok || user.Email == "" {
		err = errors.New("Invalid user object")
		return
	}

	id = user.Email
	return
}

func orderFilterFromRequest(req *http.Request) (filter core.OrderFilter, err error) {
	if err = req.ParseForm(); err != nil {
		return
	}

	queryValues := req.URL.Query()
	decoder := schema.NewDecoder()
	if err = decoder.Decode(&filter, queryValues); err != nil {
		return
	}

	filter.UserID, err = userId(req)
	return
}
