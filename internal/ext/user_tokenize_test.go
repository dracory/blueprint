package ext

import (
	"context"
	"fmt"
	"testing"

	"github.com/dracory/userstore"
	"github.com/dracory/vaultstore"
	"project/internal/testutils"
)

const (
	seedUserID   = "user-tokenize-test-user"
	vaultKeyTest = "test-key"
)

func TestUserTokenizeSuccess(t *testing.T) {
	app := testutils.Setup(
		testutils.WithUserStore(true),
		testutils.WithVaultStore(true),
	)

	ctx := context.Background()
	app.GetConfig().SetVaultStoreKey(vaultKeyTest)

	userStore := app.GetUserStore()
	user, err := testutils.SeedUser(userStore, seedUserID)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	user.SetFirstName("")
	user.SetLastName("")
	user.SetEmail("")
	user.SetPhone("")
	user.SetBusinessName("")

	firstToken, lastToken, emailToken, phoneToken, businessToken, err := UserTokenize(
		ctx,
		app.GetVaultStore(),
		vaultKeyTest,
		user,
		"John",
		"Doe",
		"john@example.com",
		"+44111222333",
		"JD Consulting",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertTokenStored(t, app.GetVaultStore(), vaultKeyTest, firstToken, "John")
	assertTokenStored(t, app.GetVaultStore(), vaultKeyTest, lastToken, "Doe")
	assertTokenStored(t, app.GetVaultStore(), vaultKeyTest, emailToken, "john@example.com")
	assertTokenStored(t, app.GetVaultStore(), vaultKeyTest, phoneToken, "+44111222333")
	assertTokenStored(t, app.GetVaultStore(), vaultKeyTest, businessToken, "JD Consulting")
}

func TestUserTokenizeReturnsErrorWhenVaultStoreNil(t *testing.T) {
	if _, _, _, _, _, err := UserTokenize(context.Background(), nil, vaultKeyTest, userstore.NewUser(), "", "", "", "", ""); err == nil {
		t.Fatalf("expected error when vault store is nil")
	}
}

func TestUserTokenizeReturnsErrorWhenUserNil(t *testing.T) {
	app := testutils.Setup(
		testutils.WithUserStore(true),
		testutils.WithVaultStore(true),
	)

	if _, _, _, _, _, err := UserTokenize(context.Background(), app.GetVaultStore(), vaultKeyTest, nil, "", "", "", "", ""); err == nil {
		t.Fatalf("expected error when user is nil")
	}
}

func TestUserTokenizePropagatesVaultErrors(t *testing.T) {
	app := testutils.Setup(
		testutils.WithUserStore(true),
		testutils.WithVaultStore(true),
	)
	ctx := context.Background()

	store := app.GetVaultStore()
	user := userstore.NewUser()

	if _, err := app.GetDB().ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", store.GetVaultTableName())); err != nil {
		t.Fatalf("failed to drop vault table: %v", err)
	}

	if _, _, _, _, _, err := UserTokenize(ctx, store, vaultKeyTest, user, "John", "Doe", "john@example.com", "", ""); err == nil {
		t.Fatalf("expected error when token upsert fails")
	}
}

func assertTokenStored(t *testing.T, store vaultstore.StoreInterface, vaultKey string, token string, expected string) {
	t.Helper()

	if token == "" {
		t.Fatalf("expected token for %s, got empty string", expected)
	}

	value, err := store.TokenRead(context.Background(), token, vaultKey)
	if err != nil {
		t.Fatalf("failed to read token %s: %v", token, err)
	}

	if value != expected {
		t.Fatalf("expected value %q, got %q", expected, value)
	}
}
