package helpers

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/userstore"
)

// UserActiveRoleHandles returns the set of handles of the active roles assigned to the user.
func UserActiveRoleHandles(ctx context.Context, application app.AppInterface, userID string) (map[string]bool, error) {
	if application == nil || application.IsDisabledUserStore() {
		return nil, errors.New("user store is not initialized")
	}

	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	roles, err := application.GetUserStore().UserRoles(ctx, userID)

	if err != nil {
		return nil, err
	}

	handles := map[string]bool{}
	for _, role := range roles {
		if role.IsActive() {
			handles[role.GetHandle()] = true
		}
	}

	return handles, nil
}

// UserHasActiveRole returns true if the user has an active role with the given handle.
func UserHasActiveRole(ctx context.Context, application app.AppInterface, user userstore.UserInterface, handle string) bool {
	if user == nil {
		return false
	}

	handles, err := UserActiveRoleHandles(ctx, application, user.GetID())

	if err != nil {
		return false
	}

	return handles[handle]
}

// UserHasAnyActiveRole returns true if the user has at least one active role with the given handles.
func UserHasAnyActiveRole(ctx context.Context, application app.AppInterface, user userstore.UserInterface, handles ...string) bool {
	if user == nil {
		return false
	}

	userHandles, err := UserActiveRoleHandles(ctx, application, user.GetID())

	if err != nil {
		return false
	}

	for _, handle := range handles {
		if userHandles[handle] {
			return true
		}
	}

	return false
}

// UserActiveRoleAssign assigns the role with the given handle to the user.
// The role must be active. The assignment is idempotent: an existing user
// role is returned unchanged.
func UserActiveRoleAssign(ctx context.Context, application app.AppInterface, userID string, handle string) error {
	if application == nil || application.IsDisabledUserStore() {
		return errors.New("user store is not initialized")
	}

	if userID == "" {
		return errors.New("user ID is empty")
	}

	role, err := application.GetUserStore().RoleFindByHandle(ctx, handle)

	if err != nil {
		return err
	}

	if role == nil {
		return errors.New("role not found: " + handle)
	}

	if !role.IsActive() {
		return errors.New("role is not active: " + handle)
	}

	_, err = application.GetUserStore().UserRoleFindByUserIDAndRoleIDOrCreate(ctx, userID, role.GetID())

	return err
}
