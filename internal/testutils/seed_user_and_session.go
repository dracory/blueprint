package testutils

import (
	"errors"
	"net/http"

	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
)

func SeedUserAndSession(userStore userstore.StoreInterface, sessionStore sessionstore.StoreInterface, userID string, r *http.Request, expiresSeconds int) (user userstore.UserInterface, session sessionstore.SessionInterface, err error) {
	if r == nil {
		return nil, nil, errors.New("request should not be nil")
	}

	user, err = SeedUser(userStore, userID)

	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, errors.New("user should not be nil")
	}

	session, err = SeedSession(sessionStore, r, user, expiresSeconds)

	if err != nil {
		return nil, nil, err
	}

	if session == nil {
		return nil, nil, errors.New("session should not be nil")
	}

	return user, session, nil
}
