package authadapters

import (
	"context"
	"errors"
	"log/slog"

	"project/internal/app"
	"project/internal/helpers"

	"github.com/dracory/blindindexstore"
	"github.com/dracory/userstore"
)

// UserFindByEmailOrCreate finds or creates a user by email — vault and
// blind-index aware. Ported verbatim from shared/session_login.go so the
// package's PostAuthHooks can drive it.
func UserFindByEmailOrCreate(application app.AppInterface, ctx context.Context, email string, status string) (userstore.UserInterface, error) {
	if application.IsDisabledUserStore() {
		return nil, errors.New("user store is nil")
	}

	if application.GetConfig().GetUserStoreVaultEnabled() {
		if application.IsDisabledVaultStore() {
			return nil, errors.New(`vault store is nil`)
		}

		userID, err := findUserIDInBlindIndex(application, ctx, email)
		if err != nil {
			return nil, err
		}

		if userID == "" {
			return UserCreate(application, ctx, email, status)
		}

		user, err := application.GetUserStore().UserFindByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			application.GetLogger().Warn("At authadapters UserFindByEmailOrCreate",
				slog.String("error", "email in blind index but user record missing"),
				slog.String("user", userID))
			return nil, nil
		}
		return user, nil
	}

	user, err := application.GetUserStore().UserFindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return UserCreate(application, ctx, email, status)
	}
	return user, nil
}

// UserFindByEmail looks up an existing user by email without creating one —
// vault-aware; prefers the record with a password set when duplicates
// exist. Returns (nil, nil) when no user exists.
func UserFindByEmail(application app.AppInterface, ctx context.Context, email string) (userstore.UserInterface, error) {
	if application.IsDisabledUserStore() {
		return nil, errors.New("user store is nil")
	}

	if application.GetConfig().GetUserStoreVaultEnabled() {
		if application.IsDisabledVaultStore() {
			return nil, errors.New(`vault store is nil`)
		}
		recordsFound, err := application.GetBlindIndexStoreEmail().SearchValueList(ctx, blindindexstore.NewSearchValueQuery().
			SetSearchValue(email).
			SetSearchType(blindindexstore.SEARCH_TYPE_EQUALS))
		if err != nil {
			return nil, err
		}
		return findUserWithPassword(application, ctx, recordsFound)
	}

	users, err := application.GetUserStore().UserList(ctx, userstore.NewUserQuery().
		SetEmail(email).
		SetLimit(100))
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	for _, user := range users {
		if user.GetPassword() != "" {
			return user, nil
		}
	}
	return users[0], nil
}

// findUserWithPassword resolves blind index hits to user records, preferring
// the first with a password set.
func findUserWithPassword(application app.AppInterface, ctx context.Context, recordsFound []blindindexstore.SearchValueInterface) (userstore.UserInterface, error) {
	var fallbackUser userstore.UserInterface
	for _, record := range recordsFound {
		user, err := application.GetUserStore().UserFindByID(ctx, record.SourceReferenceID())
		if err != nil {
			return nil, err
		}
		if user == nil {
			continue
		}
		if fallbackUser == nil {
			fallbackUser = user
		}
		if user.GetPassword() != "" {
			return user, nil
		}
	}
	return fallbackUser, nil
}

// UserCreate creates a new user with privacy-first email encryption and
// blind indexing.
func UserCreate(application app.AppInterface, ctx context.Context, email string, status string) (userstore.UserInterface, error) {
	user := userstore.NewUser().
		SetStatus(status).
		SetEmail(email)

	if application.IsDisabledUserStore() {
		return nil, errors.New("user store is nil")
	}
	if application.GetConfig().GetUserStoreVaultEnabled() && application.IsDisabledVaultStore() {
		return nil, errors.New(`vault store is nil`)
	}

	if err := application.GetUserStore().UserCreate(ctx, user); err != nil {
		return nil, err
	}
	if err := helpers.UserActiveRoleAssign(ctx, application, user.GetID(), userstore.USER_ROLE_USER); err != nil {
		return nil, err
	}
	if !application.GetConfig().GetUserStoreVaultEnabled() {
		return user, nil
	}

	emailToken, err := application.GetVaultStore().TokenCreate(ctx, email, application.GetConfig().GetVaultStoreKey(), 20)
	if err != nil {
		return nil, err
	}
	user.SetEmail(emailToken)
	if err := application.GetUserStore().UserUpdate(ctx, user); err != nil {
		return nil, err
	}

	searchValue := blindindexstore.NewSearchValue().
		SetSourceReferenceID(user.GetID()).
		SetSearchValue(email)
	if err := application.GetBlindIndexStoreEmail().SearchValueCreate(ctx, searchValue); err != nil {
		return nil, err
	}
	return user, nil
}

func findUserIDInBlindIndex(application app.AppInterface, ctx context.Context, email string) (userID string, err error) {
	recordsFound, err := application.GetBlindIndexStoreEmail().SearchValueList(ctx, blindindexstore.NewSearchValueQuery().
		SetSearchValue(email).
		SetSearchType(blindindexstore.SEARCH_TYPE_EQUALS))
	if err != nil {
		return "", err
	}
	if len(recordsFound) < 1 {
		return "", nil
	}
	return recordsFound[0].SourceReferenceID(), nil
}
