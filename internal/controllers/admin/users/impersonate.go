package admin

import (
	"errors"
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/sessionstore"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/auth"
)

func Impersonate(ss sessionstore.StoreInterface, w http.ResponseWriter, r *http.Request, userID string) error {
	if ss == nil {
		return errors.New("session store is nil")
	}

	session := sessionstore.NewSession().
		SetUserID(userID).
		SetUserAgent(r.UserAgent()).
		SetIPAddress(req.GetIP(r)).
		SetExpiresAt(carbon.Now(carbon.UTC).AddHours(2).ToDateTimeString(carbon.UTC))

	err := ss.SessionCreate(r.Context(), session)

	if err != nil {
		return err
	}

	auth.AuthCookieSet(w, r, session.GetKey())

	return nil
}
