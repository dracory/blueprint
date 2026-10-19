package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
	"github.com/dracory/userstore"
)

var _ migrator.MigrationInterface = (*StoreUserMigrate)(nil)

// StoreUserMigrate creates the userstore tables, seeds the default roles,
// and backfills user role assignments from the legacy `role` column on
// the user record.
type StoreUserMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreUserMigrate) Signature() string {
	return "2026_03_21_0021_store_user_migrate"
}

func (m *StoreUserMigrate) Description() string {
	return "Run user store MigrateUp to create user and role tables and seed default roles"
}

// defaultRoles maps role handles to display names.
var defaultRoles = map[string]string{
	userstore.USER_ROLE_SUPERUSER:     "Superuser",
	userstore.USER_ROLE_ADMINISTRATOR: "Administrator",
	userstore.USER_ROLE_MANAGER:       "Manager",
	userstore.USER_ROLE_USER:          "User",
}

func (m *StoreUserMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetUserStore()
	if store == nil {
		return errors.New("user store is not initialized")
	}

	ctx := context.Background()

	if err := store.MigrateUp(ctx); err != nil {
		return err
	}

	roleIDs := map[string]string{}
	for handle, name := range defaultRoles {
		role, err := store.RoleFindByHandleOrCreate(ctx, handle, userstore.ROLE_STATUS_ACTIVE)
		if err != nil {
			return err
		}

		if role.GetName() == "" {
			role.SetName(name)
			if err := store.RoleUpdate(ctx, role); err != nil {
				return err
			}
		}

		roleIDs[handle] = role.GetID()
	}

	// Backfill user_role assignments from the legacy role column.
	users, err := store.UserList(ctx, userstore.NewUserQuery())
	if err != nil {
		return err
	}

	for _, user := range users {
		roleID, ok := roleIDs[user.GetRole()]
		if !ok {
			continue
		}

		if _, err := store.UserRoleFindByUserIDAndRoleIDOrCreate(ctx, user.GetID(), roleID); err != nil {
			return err
		}
	}

	return nil
}

func (m *StoreUserMigrate) Down() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetUserStore()
	if store == nil {
		return errors.New("user store is not initialized")
	}

	return store.MigrateDown(context.Background())
}
