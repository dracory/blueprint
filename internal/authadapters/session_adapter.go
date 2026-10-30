// Package authadapters bridges admin's concrete stores (sessionstore,
// cachestore, userstore, taskstore) to the github.com/dracory/auth
// package's minimal interfaces, so the standalone package can drive
// login/session flows without importing project internals.
package authadapters

import (
	"context"
	"errors"
	"time"

	"github.com/dracory/auth/types"
	"github.com/dracory/sessionstore"
	"github.com/dromara/carbon/v2"
)

// Free-form value names the auth package writes on sessions, mapped to
// sessionstore's fixed columns.
const (
	// ValueRemember is the remember-me marker (package's
	// utils.RememberSessionMarker also uses "remember").
	ValueRemember = "remember"
	valueIP       = "ip"
	valueUA       = "user_agent"
)

// sessionAdapter wraps a sessionstore session to satisfy
// types.SessionInterface. sessionstore uses typed columns + a single
// session_value column; the package's free-form GetValue(name) maps:
// "remember" -> session_value, "ip" -> session_ip, "user_agent" ->
// session_user_agent. Unknown names map to session_value too (documented
// limitation — the store only carries one value column).
type sessionAdapter struct {
	s sessionstore.SessionInterface
}

// WrapSession adapts a sessionstore session to types.SessionInterface.
func WrapSession(s sessionstore.SessionInterface) types.SessionInterface {
	return &sessionAdapter{s: s}
}

func (a *sessionAdapter) GetID() string      { return a.s.GetID() }
func (a *sessionAdapter) SetID(id string)    { a.s.SetID(id) }
func (a *sessionAdapter) GetUserID() string  { return a.s.GetUserID() }
func (a *sessionAdapter) SetUserID(u string) { a.s.SetUserID(u) }
func (a *sessionAdapter) GetKey() string     { return a.s.GetKey() }
func (a *sessionAdapter) SetKey(k string)    { a.s.SetKey(k) }

func (a *sessionAdapter) GetExpiresAt() time.Time {
	c := a.s.GetExpiresAtCarbon()
	if c == nil {
		return time.Time{}
	}
	return c.StdTime()
}

func (a *sessionAdapter) SetExpiresAt(t time.Time) {
	a.s.SetExpiresAt(carbon.NewCarbon(t).ToDateTimeString(carbon.UTC))
}

func (a *sessionAdapter) GetValue(name string) string {
	switch name {
	case valueIP:
		return a.s.GetIPAddress()
	case valueUA:
		return a.s.GetUserAgent()
	default:
		return a.s.GetValue()
	}
}

func (a *sessionAdapter) SetValue(name string, value string) {
	switch name {
	case valueIP:
		a.s.SetIPAddress(value)
	case valueUA:
		a.s.SetUserAgent(value)
	default:
		a.s.SetValue(value)
	}
}

// SessionStoreAdapter adapts sessionstore.StoreInterface to
// types.SessionStoreInterface. sessionstore has no DeleteByKey (only
// ByID/object) and no DeleteByUserID — both are implemented via
// Find/List + delete so the package's contract holds.
type SessionStoreAdapter struct {
	Store sessionstore.StoreInterface
}

var _ types.SessionStoreInterface = (*SessionStoreAdapter)(nil)

// NewSession returns a blank sessionstore session wrapped as
// types.SessionInterface.
func (a *SessionStoreAdapter) NewSession() types.SessionInterface {
	return WrapSession(sessionstore.NewSession())
}

func (a *SessionStoreAdapter) SessionCreate(ctx context.Context, session types.SessionInterface) error {
	w, ok := session.(*sessionAdapter)
	if !ok {
		return errors.New("authadapters: session is not a sessionstore session")
	}
	return a.Store.SessionCreate(ctx, w.s)
}

func (a *SessionStoreAdapter) SessionFindByKey(ctx context.Context, key string) (types.SessionInterface, error) {
	s, err := a.Store.SessionFindByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}
	return WrapSession(s), nil
}

// SessionUpdate persists a mutated session — satisfies the package's
// optional middlewares.SessionUpdater interface so remember-me rotation can
// shorten a consumed session to its grace window in place.
func (a *SessionStoreAdapter) SessionUpdate(ctx context.Context, session types.SessionInterface) error {
	w, ok := session.(*sessionAdapter)
	if !ok {
		return errors.New("authadapters: session is not a sessionstore session")
	}
	return a.Store.SessionUpdate(ctx, w.s)
}

func (a *SessionStoreAdapter) SessionDeleteByKey(ctx context.Context, key string) error {
	s, err := a.Store.SessionFindByKey(ctx, key)
	if err != nil || s == nil {
		return err
	}
	return a.Store.SessionDelete(ctx, s)
}

func (a *SessionStoreAdapter) SessionDeleteByUserID(ctx context.Context, userID string) error {
	sessions, err := a.Store.SessionList(ctx, sessionstore.NewSessionQuery().SetUserID(userID))
	if err != nil {
		return err
	}
	for _, s := range sessions {
		if err := a.Store.SessionDelete(ctx, s); err != nil {
			return err
		}
	}
	return nil
}
