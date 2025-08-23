package testutils

import (
	"errors"
	"net/http"

	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/sessionstore"
	"github.com/gouniverse/userstore"
	"github.com/gouniverse/utils"
)

func SeedSession(sessionStore sessionstore.StoreInterface, r *http.Request, user userstore.UserInterface, expiresSeconds int) (sessionstore.SessionInterface, error) {
	if sessionStore == nil {
		return nil, errors.New("session store is nil")
	}

	session := sessionstore.NewSession().
		SetUserID(user.ID()).
		SetUserAgent(r.UserAgent()).
		SetIPAddress(utils.IP(r)).
		SetExpiresAt(carbon.Now(carbon.UTC).AddSeconds(expiresSeconds).ToDateTimeString(carbon.UTC))

	err := sessionStore.SessionCreate(r.Context(), session)

	if err != nil {
		return nil, err
	}

	return session, nil
}
