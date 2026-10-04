package usersonce

import (
	"context"
	"errors"
	"fmt"
	"project/internal/app"

	"github.com/dracory/base/cfmt"
	"github.com/dracory/neat/database/migrator"
	"github.com/dracory/userstore"
)

var _ migrator.MigrationInterface = (*UsersSeed)(nil)

// UsersSeed creates the default user accounts (an administrator and a
// regular user). It runs exactly once — applied seeds are recorded in the
// shared migration_tracker table under the seed_ signature, and the
// migrator skips them on subsequent boots. It is also idempotent by
// design: each account is created via find-or-create on email.
type UsersSeed struct {
	migrator.BaseMigration
	app app.AppInterface
}

// NewUsersSeed creates a new UsersSeed instance.
func NewUsersSeed(a app.AppInterface) *UsersSeed {
	return &UsersSeed{app: a}
}

func (s *UsersSeed) Signature() string {
	return "seed_2026_10_04_0002_users_default"
}

func (s *UsersSeed) Description() string {
	return "Create the default administrator and user accounts"
}

func (s *UsersSeed) Up() error {
	return seedUsers(context.Background(), s.app)
}

// Down is a no-op: seeded accounts may have drifted since the seed ran,
// so they are left in place.
func (s *UsersSeed) Down() error {
	return nil
}

func seedUsers(ctx context.Context, app app.AppInterface) error {
	if app == nil {
		return errors.New("app is nil")
	}
	if app.IsDisabledUserStore() {
		return errors.New("user store is not initialized")
	}

	store := app.GetUserStore()

	for _, u := range defaultUsers() {
		if err := findOrCreateUser(ctx, store, u); err != nil {
			return err
		}
	}

	return nil
}

// findOrCreateUser creates the account only when no user with that email
// exists — the find-or-create that makes the seed safe to re-run.
func findOrCreateUser(ctx context.Context, store userstore.StoreInterface, u defaultUser) error {
	existing, err := store.UserFindByEmail(ctx, u.email)
	if err != nil {
		return fmt.Errorf("find user %s: %w", u.email, err)
	}

	if existing == nil {
		user := userstore.NewUser().
			SetEmail(u.email).
			SetFirstName(u.firstName).
			SetLastName(u.lastName).
			SetStatus(userstore.USER_STATUS_ACTIVE)

		if err := user.SetPasswordAndHash(u.password); err != nil {
			return fmt.Errorf("hash password for %s: %w", u.email, err)
		}

		if err := store.UserCreate(ctx, user); err != nil {
			return fmt.Errorf("create user %s: %w", u.email, err)
		}
		existing = user

		cfmt.Successln("Seeded user:", u.email)
	}

	role, err := store.RoleFindByHandleOrCreate(ctx, u.role, userstore.ROLE_STATUS_ACTIVE)
	if err != nil {
		return fmt.Errorf("find role %s: %w", u.role, err)
	}

	if _, err := store.UserRoleFindByUserIDAndRoleIDOrCreate(ctx, existing.GetID(), role.GetID()); err != nil {
		return fmt.Errorf("assign role %s to %s: %w", u.role, u.email, err)
	}

	return nil
}
