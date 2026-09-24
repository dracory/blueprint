package testutils

import (
	"context"
	"errors"

	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

// SeedUser find existing or generates a new user with the given ID (no email, no names)
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

	err = userStore.UserCreate(context.Background(), user)
	if err != nil {
		return nil, err
	}

	roleHandle := ""
	if userID == test.USER_01 {
		roleHandle = userstore.USER_ROLE_USER
	}
	if userID == test.ADMIN_01 {
		roleHandle = userstore.USER_ROLE_ADMINISTRATOR
	}

	if roleHandle != "" {
		role, err := userStore.RoleFindByHandleOrCreate(context.Background(), roleHandle, userstore.ROLE_STATUS_ACTIVE)
		if err != nil {
			return nil, err
		}
		if _, err := userStore.UserRoleFindByUserIDAndRoleIDOrCreate(context.Background(), user.GetID(), role.GetID()); err != nil {
			return nil, err
		}
	}

	return user, nil
}
