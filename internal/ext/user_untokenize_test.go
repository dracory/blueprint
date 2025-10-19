package ext

import (
	"context"
	"fmt"
	"testing"

	"github.com/dracory/userstore"
	"project/internal/testutils"
)

const (
	untokenizeVaultKey = "test-key"
)

func TestUserUntokenizeSuccess(t *testing.T) {
	app := testutils.Setup(
		testutils.WithVaultStore(true),
	)

	ctx := context.Background()
	app.GetConfig().SetVaultStoreKey(untokenizeVaultKey)

	user := userstore.NewUser()

	firstToken, lastToken, emailToken, phoneToken, businessToken, err := UserTokenize(
		ctx,
		app.GetVaultStore(),
		untokenizeVaultKey,
		user,
		"John",
		"Doe",
		"john@example.com",
		"+44111222333",
		"JD Consulting",
	)
	if err != nil {
		t.Fatalf("UserTokenize failed: %v", err)
	}

	user.SetFirstName(firstToken)
	user.SetLastName(lastToken)
	user.SetEmail(emailToken)
	user.SetPhone(phoneToken)
	user.SetBusinessName(businessToken)

	firstName, lastName, email, businessName, phone, err := UserUntokenize(ctx, app, untokenizeVaultKey, user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if firstName != "John" {
		t.Fatalf("expected first name 'John', got %q", firstName)
	}

	if lastName != "Doe" {
		t.Fatalf("expected last name 'Doe', got %q", lastName)
	}

	if email != "john@example.com" {
		t.Fatalf("expected email 'john@example.com', got %q", email)
	}

	if businessName != "JD Consulting" {
		t.Fatalf("expected business name 'JD Consulting', got %q", businessName)
	}

	if phone != "+44111222333" {
		t.Fatalf("expected phone '+44111222333', got %q", phone)
	}
}

func TestUserUntokenizeReturnsErrorWhenVaultStoreNil(t *testing.T) {
	app := testutils.Setup()

	if _, _, _, _, _, err := UserUntokenize(context.Background(), app, untokenizeVaultKey, userstore.NewUser()); err == nil {
		t.Fatalf("expected error when vault store is nil")
	}
}

func TestUserUntokenizeReturnsErrorWhenUserNil(t *testing.T) {
	app := testutils.Setup(
		testutils.WithVaultStore(true),
	)

	if _, _, _, _, _, err := UserUntokenize(context.Background(), app, untokenizeVaultKey, nil); err == nil {
		t.Fatalf("expected error when user is nil")
	}
}

func TestUserUntokenizePropagatesVaultErrors(t *testing.T) {
	app := testutils.Setup(
		testutils.WithVaultStore(true),
	)

	ctx := context.Background()
	app.GetConfig().SetVaultStoreKey(untokenizeVaultKey)

	user := userstore.NewUser()

	firstToken, lastToken, emailToken, phoneToken, businessToken, err := UserTokenize(
		ctx,
		app.GetVaultStore(),
		untokenizeVaultKey,
		user,
		"John",
		"Doe",
		"john@example.com",
		"+44111222333",
		"JD Consulting",
	)
	if err != nil {
		t.Fatalf("UserTokenize failed: %v", err)
	}

	user.SetFirstName(firstToken)
	user.SetLastName(lastToken)
	user.SetEmail(emailToken)
	user.SetPhone(phoneToken)
	user.SetBusinessName(businessToken)

	if _, err := app.GetDB().ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", app.GetVaultStore().GetVaultTableName())); err != nil {
		t.Fatalf("failed to drop vault table: %v", err)
	}

	if _, _, _, _, _, err := UserUntokenize(ctx, app, untokenizeVaultKey, user); err == nil {
		t.Fatalf("expected error when vault table is missing")
	}
}
