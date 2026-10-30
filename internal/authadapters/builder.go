package authadapters

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"project/internal/app"
	"project/internal/config"
	"project/internal/emails"
	"project/internal/helpers"
	"project/internal/links"
	authrules "project/internal/rules/auth"
	"project/internal/tasks/constants"

	"github.com/dracory/auth"
	"github.com/dracory/auth/types"
	authutils "github.com/dracory/auth/utils"
	"github.com/dracory/taskstore"
	"github.com/dracory/userstore"
)

// BuildPasswordlessAuth constructs the package's passwordless auth instance
// wired to admin's stores via the adapters — the single AuthSharedInterface
// that the multi-method mounts and PostAuthHooks hang off.
func BuildPasswordlessAuth(application app.AppInterface) (types.AuthPasswordlessInterface, error) {
	tempKeys := &TempKeyAdapter{Store: application.GetCacheStore()}
	sessions := &SessionStoreAdapter{Store: application.GetSessionStore()}

	cfg := types.ConfigPasswordless{
		ConfigShared: types.ConfigShared{
			EnableRegistration:   application.GetConfig().GetRegistrationEnabled(),
			Endpoint:             "/auth",
			UrlRedirectOnSuccess: links.User().Home(),
			FuncLayout:           NewAuthLayoutFunc(application),
			UseCookies:           true,
			Logger:               application.GetLogger(),
			FuncTemporaryKeyGet:  tempKeys.Get,
			FuncTemporaryKeySet:  tempKeys.Set,
			// Session keys double as auth tokens: lookup resolves the
			// session row to its user ID.
			FuncUserFindByAuthToken: func(ctx context.Context, sessionID string, options types.UserAuthOptions) (string, error) {
				s, err := sessions.SessionFindByKey(ctx, sessionID)
				if err != nil || s == nil {
					return "", err
				}
				return s.GetUserID(), nil
			},
			// Token persistence is owned by the session store — the hook
			// exists for the interface but the session path supersedes it.
			FuncUserStoreAuthToken: func(ctx context.Context, sessionID string, userID string, options types.UserAuthOptions) error {
				return nil
			},
			FuncUserLogout: func(ctx context.Context, userID string, options types.UserAuthOptions) error {
				return sessions.SessionDeleteByUserID(ctx, userID)
			},
			CookieConfig: &types.CookieConfig{
				Name: auth.CookieName,
				Path: "/",
				// HTTP envs must not set Secure or the browser drops the
				// auth cookie entirely — local envs may be plain HTTP.
				Secure: secureFlag(application),
			},
		},
		FuncUserFindByEmail: func(ctx context.Context, email string, options types.UserAuthOptions) (string, error) {
			user, err := UserFindByEmail(application, ctx, email)
			if err != nil {
				application.GetLogger().Error("authadapters FuncUserFindByEmail", slog.String("error", err.Error()), slog.String("email", email))
				return "", err
			}
			if user == nil {
				return "", nil
			}
			return user.GetID(), nil
		},
		FuncEmailSend: func(ctx context.Context, email string, subject string, body string) error {
			err := FuncEmailSend(application)(ctx, email, subject, body)
			if err != nil {
				application.GetLogger().Error("authadapters FuncEmailSend", slog.String("error", err.Error()))
			}
			return err
		},
		// Brand the OTP email — the package's default template is plain.
		// A panic inside the template must not kill the send (DeliverSecret
		// has no recover) — log it and fall back to a minimal body.
		FuncEmailTemplateLoginCode: func(ctx context.Context, email string, loginLink string, options types.UserAuthOptions) (body string) {
			defer func() {
				if r := recover(); r != nil {
					application.GetLogger().Error("authadapters FuncEmailTemplateLoginCode panic",
						slog.Any("panic", r))
					body = "Your login code: " + loginLink
				}
			}()
			otp := emails.NewEmailOtp(application)
			return otp.BodyTemplate(otp.Subject(), loginLink)
		},
		FuncUserRegister: func(ctx context.Context, email string, firstName string, lastName string, options types.UserAuthOptions) error {
			user, err := UserCreate(application, ctx, email, userstore.USER_STATUS_ACTIVE)
			if err != nil {
				application.GetLogger().Error("authadapters FuncUserRegister", slog.String("error", err.Error()), slog.String("email", email))
				return err
			}
			if user == nil {
				return nil
			}
			return nil
		},
	}

	a, err := auth.NewPasswordlessAuth(cfg)
	if err != nil {
		return nil, err
	}

	if impl, ok := a.(sessionSetter); ok {
		impl.SetSessionStore(sessions)
		impl.SetFuncTemporaryKeyDelete(tempKeys.Delete)
		impl.SetPostAuthHooks(postAuthHooks(application, sessions))
		impl.SetFuncEmailEnqueue(emailEnqueue(application))
		rememberMeSettings(application, impl)
	}
	return a, nil
}

// sessionSetter captures the optional-hook surface the auth
// implementation exposes via type assertion (adding these to
// AuthSharedInterface would break external implementers).
type sessionSetter interface {
	SetSessionStore(types.SessionStoreInterface)
	SetFuncTemporaryKeyDelete(func(key string) error)
	SetPostAuthHooks(types.PostAuthHooks)
	SetFuncEmailEnqueue(func(ctx context.Context, deliveryKey string) error)
	SetRememberMeEnabled(bool)
	SetRememberMeDays(int)
	SetRememberCookieName(string)
}

// rememberMeSettings pushes the app's remember-me config onto a built
// package auth instance — the package pipeline then issues a marked
// remember session + cookie when the user opts in on the login page.
func rememberMeSettings(application app.AppInterface, impl sessionSetter) {
	impl.SetRememberMeEnabled(application.GetConfig().GetRememberMeEnabled())
	impl.SetRememberMeDays(application.GetConfig().GetRememberMeDays())
	impl.SetRememberCookieName(config.COOKIE_NAME_REMEMBER_TOKEN)
}

// emailEnqueue returns the package's queue-safe enqueue hook: it pushes the
// opaque delivery key onto AuthEmailDeliveryTask so email credentials never
// enter the task queue. Returns nil when the task store is disabled — the
// package then falls back to the synchronous FuncEmailSend hook, so hosts
// without a queue still deliver.
func emailEnqueue(application app.AppInterface) func(ctx context.Context, deliveryKey string) error {
	if application.IsDisabledTaskStore() {
		return nil
	}
	return func(ctx context.Context, deliveryKey string) error {
		_, err := application.GetTaskStore().TaskDefinitionEnqueueByAlias(
			ctx,
			taskstore.DefaultQueueName,
			constants.EmailAuthDeliveryTaskAlias,
			map[string]any{"delivery_key": deliveryKey},
		)
		return err
	}
}

// secureFlag returns the Secure cookie flag for the environment —
// production-only, mirroring shared.CreateSession.
func secureFlag(application app.AppInterface) string {
	if application.GetConfig().IsEnvDevelopment() || application.GetConfig().IsEnvLocal() {
		return "no"
	}
	return "yes"
}

// postAuthHooks wires the package's shared post-auth pipeline:
// allowlist → find-or-create (vault/blind-index) → active → session →
// cookie → redirect + registration rules.
func postAuthHooks(application app.AppInterface, sessions *SessionStoreAdapter) types.PostAuthHooks {
	return types.PostAuthHooks{
		EmailAllowed: func(ctx context.Context, email string) (bool, string) {
			rule := authrules.NewEmailAllowedRule(application, email)
			return rule.Fails() == false, rule.FailMessageFirst()
		},
		UserFindByEmailOrCreate: func(ctx context.Context, email string, options types.UserAuthOptions) (string, error) {
			user, err := UserFindByEmailOrCreate(application, ctx, email, userstore.USER_STATUS_ACTIVE)
			if err != nil || user == nil {
				return "", err
			}
			return user.GetID(), nil
		},
		UserIsActive:  UserIsActiveHook(application),
		SessionCreate: IssueAuthSession(application, sessions),
		RedirectURL: func(ctx context.Context, userID string) string {
			// Mirrors shared.calculateRedirectURL: admins → admin panel,
			// incomplete registration → register page, else user home.
			user, err := application.GetUserStore().UserFindByID(ctx, userID)
			if err != nil || user == nil {
				return links.User().Home()
			}
			redirectURL := links.User().Home()
			if helpers.UserHasAnyActiveRole(ctx, application, user,
				userstore.USER_ROLE_MANAGER,
				userstore.USER_ROLE_ADMINISTRATOR,
				userstore.USER_ROLE_SUPERUSER) {
				redirectURL = links.Admin().Home()
			}
			if !user.IsRegistrationCompleted() {
				redirectURL = helpers.ToFlashInfoURL(application.GetCacheStore(),
					"Thank you for logging in. Please complete your data to finish your registration",
					links.Auth().Register(), 5)
			}
			return redirectURL
		},
		SetAuthCookie: func(w http.ResponseWriter, r *http.Request, token string) {
			opts := []types.CookieOption{}
			if application.GetConfig().IsEnvDevelopment() || application.GetConfig().IsEnvLocal() {
				opts = append(opts, types.WithSecure(false))
			}
			auth.AuthCookieSet(w, r, token, opts...)
		},
		IssueRememberSession: IssueRememberSession(application, sessions),
	}
}

// UserIsActiveHook rejects inactive/banned accounts — shared by the
// post-auth pipeline and the remember-me middleware's revocation check.
func UserIsActiveHook(application app.AppInterface) func(ctx context.Context, userID string, options types.UserAuthOptions) bool {
	return func(ctx context.Context, userID string, options types.UserAuthOptions) bool {
		user, err := application.GetUserStore().UserFindByID(ctx, userID)
		if err != nil || user == nil {
			return false
		}
		return authrules.NewUserActiveRule(user).Fails() == false
	}
}

// IssueAuthSession mints the short-lived auth session (2h production, 4h
// dev) — used by the post-auth pipeline and by the remember-me middleware
// when it exchanges a remember token for a fresh session.
func IssueAuthSession(application app.AppInterface, sessions *SessionStoreAdapter) func(ctx context.Context, userID string, options types.UserAuthOptions) (types.SessionInterface, error) {
	return func(ctx context.Context, userID string, options types.UserAuthOptions) (types.SessionInterface, error) {
		ttl := 2 * time.Hour // production default; dev extends below
		if application.GetConfig().IsEnvDevelopment() || application.GetConfig().IsEnvLocal() {
			ttl = 4 * time.Hour
		}
		s := sessions.NewSession()
		s.SetUserID(userID)
		s.SetExpiresAt(time.Now().Add(ttl))
		s.SetValue("ip", options.UserIp)
		s.SetValue("user_agent", options.UserAgent)
		if err := sessions.SessionCreate(ctx, s); err != nil {
			return nil, err
		}
		return s, nil
	}
}

// IssueRememberSession creates the marked long-lived session behind the
// remember cookie (REMEMBER_ME_DAYS, default 30). The "remember" marker
// lands in session_value via the adapter — the middleware only accepts
// marked sessions, and a marked key can never act as an auth credential.
func IssueRememberSession(application app.AppInterface, sessions *SessionStoreAdapter) func(ctx context.Context, userID string, options types.UserAuthOptions) (types.SessionInterface, error) {
	return func(ctx context.Context, userID string, options types.UserAuthOptions) (types.SessionInterface, error) {
		if application.GetConfig() == nil || !application.GetConfig().GetRememberMeEnabled() {
			return nil, nil
		}
		days := application.GetConfig().GetRememberMeDays()
		if days < 1 {
			days = 30
		}
		s := sessions.NewSession()
		s.SetUserID(userID)
		s.SetExpiresAt(time.Now().Add(time.Duration(days) * 24 * time.Hour))
		s.SetValue("ip", options.UserIp)
		s.SetValue("user_agent", options.UserAgent)
		s.SetValue(authutils.RememberSessionMarker, "remember")
		if err := sessions.SessionCreate(ctx, s); err != nil {
			return nil, err
		}
		return s, nil
	}
}

// BuildPasswordAuth constructs the package's username(email)+password auth
// instance wired to admin's stores — parallel to BuildPasswordlessAuth so
// multi-method mounts can mix both impls.
func BuildPasswordAuth(application app.AppInterface) (types.AuthPasswordInterface, error) {
	tempKeys := &TempKeyAdapter{Store: application.GetCacheStore()}
	sessions := &SessionStoreAdapter{Store: application.GetSessionStore()}

	cfg := types.ConfigUsernameAndPassword{
		ConfigShared: types.ConfigShared{
			EnableRegistration:   application.GetConfig().GetRegistrationEnabled(),
			Endpoint:             "/auth",
			UrlRedirectOnSuccess: links.User().Home(),
			FuncLayout:           NewAuthLayoutFunc(application),
			UseCookies:           true,
			Logger:               application.GetLogger(),
			FuncTemporaryKeyGet:  tempKeys.Get,
			FuncTemporaryKeySet:  tempKeys.Set,
			FuncUserFindByAuthToken: func(ctx context.Context, sessionID string, options types.UserAuthOptions) (string, error) {
				s, err := sessions.SessionFindByKey(ctx, sessionID)
				if err != nil || s == nil {
					return "", err
				}
				return s.GetUserID(), nil
			},
			FuncUserStoreAuthToken: func(ctx context.Context, sessionID string, userID string, options types.UserAuthOptions) error {
				return nil
			},
			FuncUserLogout: func(ctx context.Context, userID string, options types.UserAuthOptions) error {
				return sessions.SessionDeleteByUserID(ctx, userID)
			},
			CookieConfig: &types.CookieConfig{
				Name:   auth.CookieName,
				Path:   "/",
				Secure: secureFlag(application),
			},
		},
		LabelUsername: "Email",
		// Admin uses email as the username everywhere — the package's
		// "username" hooks map onto the email column/vault lookups.
		FuncUserFindByUsername: func(ctx context.Context, username string, firstName string, lastName string, options types.UserAuthOptions) (string, error) {
			user, err := UserFindByEmail(application, ctx, username)
			if err != nil {
				application.GetLogger().Error("authadapters FuncUserFindByUsername", slog.String("error", err.Error()))
				return "", err
			}
			if user == nil {
				return "", nil
			}
			return user.GetID(), nil
		},
		FuncUserLogin: func(ctx context.Context, username string, password string, options types.UserAuthOptions) (string, error) {
			user, err := UserFindByEmail(application, ctx, username)
			if err != nil {
				application.GetLogger().Error("authadapters FuncUserLogin", slog.String("error", err.Error()))
				return "", err
			}
			if user == nil || !user.PasswordCompare(password) {
				return "", errors.New("invalid credentials")
			}
			return user.GetID(), nil
		},
		FuncUserRegister: func(ctx context.Context, username string, password string, firstName string, lastName string, options types.UserAuthOptions) error {
			user, err := UserCreate(application, ctx, username, userstore.USER_STATUS_ACTIVE)
			if err != nil {
				application.GetLogger().Error("authadapters FuncUserRegister", slog.String("error", err.Error()))
				return err
			}
			if user == nil {
				return errors.New("user was not created")
			}
			if err := user.SetPasswordAndHash(password); err != nil {
				return err
			}
			return application.GetUserStore().UserUpdate(ctx, user)
		},
		FuncUserPasswordChange: func(ctx context.Context, username string, newPassword string, options types.UserAuthOptions) error {
			user, err := UserFindByEmail(application, ctx, username)
			if err != nil || user == nil {
				return err
			}
			if err := user.SetPasswordAndHash(newPassword); err != nil {
				return err
			}
			return application.GetUserStore().UserUpdate(ctx, user)
		},
		// The password impl's send hook is userID-scoped — resolve the
		// email via the user store before handing to the shared sender.
		FuncEmailSend: func(ctx context.Context, userID string, subject string, body string) error {
			user, err := application.GetUserStore().UserFindByID(ctx, userID)
			if err != nil || user == nil {
				return err
			}
			err = FuncEmailSend(application)(ctx, user.GetEmail(), subject, body)
			if err != nil {
				application.GetLogger().Error("authadapters FuncEmailSend(password)", slog.String("error", err.Error()))
			}
			return err
		},
		// Brand the reset email — package default is plain. Panic-safe for
		// the same reason as the OTP template hook.
		FuncEmailTemplatePasswordRestore: func(ctx context.Context, userID string, resetLink string, options types.UserAuthOptions) (body string) {
			defer func() {
				if r := recover(); r != nil {
					application.GetLogger().Error("authadapters FuncEmailTemplatePasswordRestore panic",
						slog.Any("panic", r))
					body = "Reset your password: " + resetLink
				}
			}()
			pr := emails.NewEmailPasswordReset(application)
			return pr.BodyTemplate(pr.Subject(), resetLink)
		},
	}

	a, err := auth.NewUsernameAndPasswordAuth(cfg)
	if err != nil {
		return nil, err
	}
	if impl, ok := a.(sessionSetter); ok {
		impl.SetSessionStore(sessions)
		impl.SetFuncTemporaryKeyDelete(tempKeys.Delete)
		impl.SetPostAuthHooks(postAuthHooks(application, sessions))
		impl.SetFuncEmailEnqueue(emailEnqueue(application))
		rememberMeSettings(application, impl)
	}
	return a, nil
}

// BuildAuthKnightAuth constructs the package's AuthKnight auth instance —
// the external service verifies the user, then the shared post-auth
// pipeline owns find-or-create, session, cookie, and redirect.
func BuildAuthKnightAuth(application app.AppInterface) (types.AuthAuthKnightInterface, error) {
	tempKeys := &TempKeyAdapter{Store: application.GetCacheStore()}
	sessions := &SessionStoreAdapter{Store: application.GetSessionStore()}
	hooks := postAuthHooks(application, sessions)

	cfg := types.ConfigAuthKnight{
		ConfigShared: types.ConfigShared{
			EnableRegistration:   application.GetConfig().GetRegistrationEnabled(),
			Endpoint:             "/auth",
			UrlRedirectOnSuccess: links.User().Home(),
			FuncLayout:           NewAuthLayoutFunc(application),
			UseCookies:           true,
			Logger:               application.GetLogger(),
			FuncTemporaryKeyGet:  tempKeys.Get,
			FuncTemporaryKeySet:  tempKeys.Set,
			FuncUserFindByAuthToken: func(ctx context.Context, sessionID string, options types.UserAuthOptions) (string, error) {
				s, err := sessions.SessionFindByKey(ctx, sessionID)
				if err != nil || s == nil {
					return "", err
				}
				return s.GetUserID(), nil
			},
			FuncUserStoreAuthToken: func(ctx context.Context, sessionID string, userID string, options types.UserAuthOptions) error {
				return nil
			},
			FuncUserLogout: func(ctx context.Context, userID string, options types.UserAuthOptions) error {
				return sessions.SessionDeleteByUserID(ctx, userID)
			},
			CookieConfig: &types.CookieConfig{
				Name:   auth.CookieName,
				Path:   "/",
				Secure: secureFlag(application),
			},
		},
		FuncUserFindByEmail: func(ctx context.Context, email string, options types.UserAuthOptions) (string, error) {
			user, err := UserFindByEmail(application, ctx, email)
			if err != nil {
				application.GetLogger().Error("authadapters FuncUserFindByEmail(authknight)", slog.String("error", err.Error()))
				return "", err
			}
			if user == nil {
				return "", nil
			}
			return user.GetID(), nil
		},
		FuncUserRegister: func(ctx context.Context, email string, firstName string, lastName string, options types.UserAuthOptions) (string, error) {
			user, err := UserCreate(application, ctx, email, userstore.USER_STATUS_ACTIVE)
			if err != nil {
				application.GetLogger().Error("authadapters FuncUserRegister(authknight)", slog.String("error", err.Error()))
				return "", err
			}
			if user == nil {
				return "", errors.New("user was not created")
			}
			return user.GetID(), nil
		},
		// Legacy-path fallback; the post-auth pipeline normally wins.
		FuncRedirectURL: hooks.RedirectURL,
	}

	a, err := auth.NewAuthKnightAuth(cfg)
	if err != nil {
		return nil, err
	}
	if impl, ok := a.(sessionSetter); ok {
		impl.SetSessionStore(sessions)
		impl.SetFuncTemporaryKeyDelete(tempKeys.Delete)
		impl.SetPostAuthHooks(hooks)
		impl.SetFuncEmailEnqueue(emailEnqueue(application))
		rememberMeSettings(application, impl)
	}
	return a, nil
}
