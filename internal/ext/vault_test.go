package ext

import (
	"context"
	"fmt"
	"testing"

	"project/internal/testutils"
)

const vaultTestKey = "test-key"

func TestVaultTokenUpsertCreatesTokenWhenMissing(t *testing.T) {
	app := testutils.Setup(
		testutils.WithVaultStore(true),
	)

	ctx := context.Background()
	store := app.GetVaultStore()

	if store == nil {
		t.Fatalf("expected vault store to be initialized")
	}

	token, err := VaultTokenUpsert(ctx, store, vaultTestKey, "", "value")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("expected generated token, got empty string")
	}

	exists, err := store.TokenExists(ctx, token)
	if err != nil {
		t.Fatalf("failed to check token existence: %v", err)
	}

	if !exists {
		t.Fatalf("expected token to exist")
	}

	stored, err := store.TokenRead(ctx, token, vaultTestKey)
	if err != nil {
		t.Fatalf("failed to read token: %v", err)
	}

	if stored != "value" {
		t.Fatalf("expected stored value 'value', got %q", stored)
	}
}

func TestVaultTokenUpsertUpdatesExistingToken(t *testing.T) {
	app := testutils.Setup(
		testutils.WithVaultStore(true),
	)

	ctx := context.Background()
	store := app.GetVaultStore()

	if store == nil {
		t.Fatalf("expected vault store to be initialized")
	}

	token, err := VaultTokenUpsert(ctx, store, vaultTestKey, "", "initial")
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	if token == "" {
		t.Fatalf("expected generated token, got empty string")
	}

	exists, err := store.TokenExists(ctx, token)
	if err != nil {
		t.Fatalf("failed to check token existence: %v", err)
	}

	if !exists {
		t.Fatalf("expected token to exist")
	}

	updatedToken, err := VaultTokenUpsert(ctx, store, vaultTestKey, token, "updated")
	if err != nil {
		t.Fatalf("expected no error updating token, got %v", err)
	}

	if updatedToken != token {
		t.Fatalf("expected token %q to remain unchanged, got %q", token, updatedToken)
	}

	exists, err = store.TokenExists(ctx, token)
	if err != nil {
		t.Fatalf("failed to check token existence: %v", err)
	}

	if !exists {
		t.Fatalf("expected token to exist")
	}

	stored, err := store.TokenRead(ctx, token, vaultTestKey)
	if err != nil {
		t.Fatalf("failed to read token: %v", err)
	}

	if stored != "updated" {
		t.Fatalf("expected stored value 'updated', got %q", stored)
	}
}

func TestVaultTokenUpsertReturnsErrorWhenCreateFails(t *testing.T) {
	app := testutils.Setup(
		testutils.WithVaultStore(true),
	)

	ctx := context.Background()
	store := app.GetVaultStore()

	if _, err := app.GetDB().ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", store.GetVaultTableName())); err != nil {
		t.Fatalf("failed to drop vault table: %v", err)
	}

	if _, err := VaultTokenUpsert(ctx, store, vaultTestKey, "", "value"); err == nil {
		t.Fatalf("expected error when token creation fails")
	}
}
