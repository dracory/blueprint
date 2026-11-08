package usersonce

import (
	"context"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/userstore"
)

func TestUsersSeedSignature(t *testing.T) {
	seed := NewUsersSeed(nil)
	if seed.Signature() != "seed_2026_10_04_0002_users_default" {
		t.Errorf("Signature() = %q", seed.Signature())
	}
	if seed.Description() == "" {
		t.Error("Description() should not be empty")
	}
}

func TestUsersSeedNilApp(t *testing.T) {
	seed := NewUsersSeed(nil)
	if err := seed.Up(); err == nil {
		t.Fatal("expected error with nil app")
	}
}

func TestUsersSeedDisabledStore(t *testing.T) {
	app := testutils.Setup() // user store disabled
	defer app.Close()

	seed := NewUsersSeed(app)
	if err := seed.Up(); err == nil {
		t.Fatal("expected error when user store is disabled")
	}
}

func TestUsersSeedCreatesDefaultUsers(t *testing.T) {
	app := testutils.Setup(testutils.WithUserStore(true))
	defer app.Close()

	seed := NewUsersSeed(app)
	if err := seed.Up(); err != nil {
		t.Fatalf("Up() failed: %v", err)
	}

	ctx := context.Background()
	store := app.GetUserStore()

	for _, u := range defaultUsers() {
		user, err := store.UserFindByEmail(ctx, u.email)
		if err != nil {
			t.Fatalf("UserFindByEmail(%s) failed: %v", u.email, err)
		}
		if user == nil {
			t.Fatalf("expected user %s to exist", u.email)
		}
		if user.GetFirstName() != u.firstName {
			t.Errorf("first name = %q, want %q", user.GetFirstName(), u.firstName)
		}

		roles, err := store.UserRoles(ctx, user.GetID())
		if err != nil {
			t.Fatalf("UserRoles failed: %v", err)
		}
		found := false
		for _, r := range roles {
			if r.GetHandle() == u.role {
				found = true
			}
		}
		if !found {
			t.Errorf("expected role %q for user %s", u.role, u.email)
		}
	}
}

func TestUsersSeedIsIdempotent(t *testing.T) {
	app := testutils.Setup(testutils.WithUserStore(true))
	defer app.Close()

	seed := NewUsersSeed(app)
	if err := seed.Up(); err != nil {
		t.Fatalf("first Up() failed: %v", err)
	}
	if err := seed.Up(); err != nil {
		t.Fatalf("second Up() failed: %v", err)
	}

	ctx := context.Background()
	users, err := app.GetUserStore().UserList(ctx, userstore.NewUserQuery())
	if err != nil {
		t.Fatalf("UserList failed: %v", err)
	}
	if len(users) != len(defaultUsers()) {
		t.Errorf("expected %d users, got %d", len(defaultUsers()), len(users))
	}
}
