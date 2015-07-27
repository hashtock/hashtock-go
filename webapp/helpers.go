package webapp

import (
	"errors"
	"net/http"

	"github.com/gorilla/context"
	authCore "github.com/hashtock/auth/core"
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
