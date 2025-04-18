package admin

import (
	"errors"
	"net/http"
	"project/config"

	"github.com/dracory/base/req"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/auth"
	"github.com/gouniverse/sessionstore"
)

func Impersonate(w http.ResponseWriter, r *http.Request, userID string) error {
	if config.SessionStore == nil {
		return errors.New("session store is nil")
	}

	session := sessionstore.NewSession().
		SetUserID(userID).
		SetUserAgent(r.UserAgent()).
		SetIPAddress(req.IP(r)).
		SetExpiresAt(carbon.Now(carbon.UTC).AddHours(2).ToDateTimeString(carbon.UTC))

	if config.IsEnvDevelopment() {
		session.SetExpiresAt(carbon.Now(carbon.UTC).AddHours(4).ToDateTimeString(carbon.UTC))
	}

	err := config.SessionStore.SessionCreate(r.Context(), session)

	if err != nil {
		config.Logger.Error("At Impersonate Error: ", "error", err.Error())
		return err
	}

	auth.AuthCookieSet(w, r, session.GetKey())

	return nil
}
