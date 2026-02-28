package ext

import (
	"context"
	"fmt"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/userstore"
)

const (
	untokenizeVaultKey = "test-key"
)

func TestUserUntokenizeSuccess(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithVaultStore(true),
	)

	ctx := context.Background()
	registry.GetConfig().SetVaultStoreKey(untokenizeVaultKey)

	user := userstore.NewUser()

	firstToken, lastToken, emailToken, phoneToken, businessToken, err := UserTokenize(
		ctx,
		registry.GetVaultStore(),
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

	firstName, lastName, email, businessName, phone, err := UserUntokenize(ctx, registry, untokenizeVaultKey, user)
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
	registry := testutils.Setup()

	if _, _, _, _, _, err := UserUntokenize(context.Background(), registry, untokenizeVaultKey, userstore.NewUser()); err == nil {
		t.Fatalf("expected error when vault store is nil")
	}
}

func TestUserUntokenizeReturnsErrorWhenUserNil(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithVaultStore(true),
	)

	if _, _, _, _, _, err := UserUntokenize(context.Background(), registry, untokenizeVaultKey, nil); err == nil {
		t.Fatalf("expected error when user is nil")
	}
}

func TestUserUntokenizePropagatesVaultErrors(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithVaultStore(true),
	)

	ctx := context.Background()
	registry.GetConfig().SetVaultStoreKey(untokenizeVaultKey)

	user := userstore.NewUser()

	firstToken, lastToken, emailToken, phoneToken, businessToken, err := UserTokenize(
		ctx,
		registry.GetVaultStore(),
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

	if _, err := registry.GetDatabase().ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", registry.GetVaultStore().GetVaultTableName())); err != nil {
		t.Fatalf("failed to drop vault table: %v", err)
	}

	if _, _, _, _, _, err := UserUntokenize(ctx, registry, untokenizeVaultKey, user); err == nil {
		t.Fatalf("expected error when vault table is missing")
	}
}

func TestUserUntokenizeTransparently_VaultDisabledReturnsPlainValues(t *testing.T) {
	registry := testutils.Setup()
	ctx := context.Background()

	user := userstore.NewUser()
	user.SetFirstName("John")
	user.SetLastName("Doe")
	user.SetEmail("john@example.com")
	user.SetPhone("+44111222333")
	user.SetBusinessName("JD Consulting")

	email, firstName, lastName, businessName, phone, err := UserUntokenizeTransparently(ctx, registry, user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if email != "john@example.com" {
		t.Fatalf("expected email 'john@example.com', got %q", email)
	}
	if firstName != "John" {
		t.Fatalf("expected first name 'John', got %q", firstName)
	}
	if lastName != "Doe" {
		t.Fatalf("expected last name 'Doe', got %q", lastName)
	}
	if businessName != "JD Consulting" {
		t.Fatalf("expected business name 'JD Consulting', got %q", businessName)
	}
	if phone != "+44111222333" {
		t.Fatalf("expected phone '+44111222333', got %q", phone)
	}
}

func TestUserUntokenizeTransparently_VaultEnabledUntokenizes(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithVaultStore(true),
	)
	ctx := context.Background()
	registry.GetConfig().SetVaultStoreKey(untokenizeVaultKey)
	registry.GetConfig().SetUserStoreVaultEnabled(true)

	user := userstore.NewUser()

	firstToken, lastToken, emailToken, phoneToken, businessToken, err := UserTokenize(
		ctx,
		registry.GetVaultStore(),
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

	email, firstName, lastName, businessName, phone, err := UserUntokenizeTransparently(ctx, registry, user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if email != "john@example.com" {
		t.Fatalf("expected email 'john@example.com', got %q", email)
	}
	if firstName != "John" {
		t.Fatalf("expected first name 'John', got %q", firstName)
	}
	if lastName != "Doe" {
		t.Fatalf("expected last name 'Doe', got %q", lastName)
	}
	if businessName != "JD Consulting" {
		t.Fatalf("expected business name 'JD Consulting', got %q", businessName)
	}
	if phone != "+44111222333" {
		t.Fatalf("expected phone '+44111222333', got %q", phone)
	}
}

func TestUserUntokenizeTransparently_ReturnsErrorWhenUserNil(t *testing.T) {
	registry := testutils.Setup()
	ctx := context.Background()

	if _, _, _, _, _, err := UserUntokenizeTransparently(ctx, registry, nil); err == nil {
		t.Fatalf("expected error when user is nil")
	}
}
