package testutils

import (
	"errors"
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/sessionstore"
	"github.com/dracory/userstore"
	"github.com/dromara/carbon/v2"
)

func SeedSession(sessionStore sessionstore.StoreInterface, r *http.Request, user userstore.UserInterface, expiresSeconds int) (sessionstore.SessionInterface, error) {
	if sessionStore == nil {
		return nil, errors.New("session store is nil")
	}

	session := sessionstore.NewSession().
		SetUserID(user.ID()).
		SetUserAgent(r.UserAgent()).
		SetIPAddress(req.GetIP(r)).
		SetExpiresAt(carbon.Now(carbon.UTC).AddSeconds(expiresSeconds).ToDateTimeString(carbon.UTC))

	err := sessionStore.SessionCreate(r.Context(), session)

	if err != nil {
		return nil, err
	}

	return session, nil
}
