package testutils

import (
	"context"
	"errors"

	"github.com/gouniverse/userstore"
)

// SeedUser find existing or generates a new user with the given ID
func SeedUser(userStore userstore.StoreInterface, userID string) (userstore.UserInterface, error) {
	if userStore == nil {
		return nil, errors.New("userstore is not configured")
	}

	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	user, err := userStore.UserFindByID(context.Background(), userID)

	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	user = userstore.NewUser().
		SetID(userID).
		SetStatus(userstore.USER_STATUS_ACTIVE)

	if userID == USER_01 {
		user.SetRole(userstore.USER_ROLE_USER)
	}

	if userID == ADMIN_01 {
		user.SetRole(userstore.USER_ROLE_ADMINISTRATOR)
	}

	err = userStore.UserCreate(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
