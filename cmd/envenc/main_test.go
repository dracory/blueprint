package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dracory/envenc"
)

// TestNewCliCreatesInstance verifies that NewCli() creates a valid CLI instance
func TestNewCliCreatesInstance(t *testing.T) {
	cli := envenc.NewCli()
	if cli == nil {
		t.Fatal("NewCli() returned nil")
	}
}

// TestEncryptDecryptRoundtrip verifies basic encryption/decryption works
func TestEncryptDecryptRoundtrip(t *testing.T) {
	password := "test-password-123"
	plaintext := "test-secret-value"

	// Encrypt
	encrypted, err := envenc.Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	if encrypted == "" {
		t.Fatal("Encrypt returned empty string")
	}
	if encrypted == plaintext {
		t.Fatal("Encrypt did not change the input")
	}

	// Decrypt
	decrypted, err := envenc.Decrypt(encrypted, password)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("Decrypt returned wrong value: got %q, want %q", decrypted, plaintext)
	}
}

// TestEncryptDecryptWithWrongPassword verifies decryption fails with wrong password
func TestEncryptDecryptWithWrongPassword(t *testing.T) {
	password := "correct-password"
	wrongPassword := "wrong-password"
	plaintext := "test-secret-value"

	encrypted, err := envenc.Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = envenc.Decrypt(encrypted, wrongPassword)
	if err == nil {
		t.Error("Decrypt with wrong password should have failed")
	}
}

// TestObfuscateDeobfuscateRoundtrip verifies obfuscation roundtrip
func TestObfuscateDeobfuscateRoundtrip(t *testing.T) {
	input := "test-value-123"

	obfuscated, err := envenc.Obfuscate(input)
	if err != nil {
		t.Fatalf("Obfuscate failed: %v", err)
	}
	if obfuscated == "" {
		t.Fatal("Obfuscate returned empty string")
	}

	deobfuscated, err := envenc.Deobfuscate(obfuscated)
	if err != nil {
		t.Fatalf("Deobfuscate failed: %v", err)
	}
	if deobfuscated != input {
		t.Errorf("Deobfuscate returned wrong value: got %q, want %q", deobfuscated, input)
	}
}

// TestVaultFileOperations verifies vault file creation and key operations
func TestVaultFileOperations(t *testing.T) {
	// Create temp directory for test vault
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test.vault")
	password := "vault-password-123"

	// Initialize vault
	err := envenc.Init(vaultPath, password)
	if err != nil {
		t.Fatalf("Init vault failed: %v", err)
	}

	// Verify vault file was created
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		t.Fatal("Vault file was not created")
	}

	// Set a key
	keyName := "TEST_KEY"
	keyValue := "test-value-456"
	err = envenc.KeySet(vaultPath, password, keyName, keyValue)
	if err != nil {
		t.Fatalf("KeySet failed: %v", err)
	}

	// Verify key exists
	exists, err := envenc.KeyExists(vaultPath, password, keyName)
	if err != nil {
		t.Fatalf("KeyExists failed: %v", err)
	}
	if !exists {
		t.Error("Key should exist after KeySet")
	}

	// Get key value
	retrievedValue, err := envenc.KeyGet(vaultPath, password, keyName)
	if err != nil {
		t.Fatalf("KeyGet failed: %v", err)
	}
	if retrievedValue != keyValue {
		t.Errorf("KeyGet returned wrong value: got %q, want %q", retrievedValue, keyValue)
	}

	// List keys
	keys, err := envenc.KeyListFromFile(vaultPath, password)
	if err != nil {
		t.Fatalf("KeyListFromFile failed: %v", err)
	}
	if len(keys) == 0 {
		t.Error("KeyListFromFile returned empty map")
	}
	if _, ok := keys[keyName]; !ok {
		t.Errorf("KeyListFromFile does not contain key %q", keyName)
	}
}

// TestKeyRemove verifies key removal works
func TestKeyRemove(t *testing.T) {
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test.vault")
	password := "vault-password-123"
	keyName := "REMOVE_KEY"
	keyValue := "value-to-remove"

	// Initialize and set key
	err := envenc.Init(vaultPath, password)
	if err != nil {
		t.Fatalf("Init vault failed: %v", err)
	}

	err = envenc.KeySet(vaultPath, password, keyName, keyValue)
	if err != nil {
		t.Fatalf("KeySet failed: %v", err)
	}

	// Remove key
	err = envenc.KeyRemove(vaultPath, password, keyName)
	if err != nil {
		t.Fatalf("KeyRemove failed: %v", err)
	}

	// Verify key no longer exists
	exists, err := envenc.KeyExists(vaultPath, password, keyName)
	if err != nil {
		t.Fatalf("KeyExists failed: %v", err)
	}
	if exists {
		t.Error("Key should not exist after KeyRemove")
	}
}

// TestHydrateEnvFromFile verifies environment hydration from vault file
func TestHydrateEnvFromFile(t *testing.T) {
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test.vault")
	password := "vault-password-123"
	keyName := "HYDRATE_TEST_KEY"
	keyValue := "hydrated-value-789"

	// Initialize and set key
	err := envenc.Init(vaultPath, password)
	if err != nil {
		t.Fatalf("Init vault failed: %v", err)
	}

	err = envenc.KeySet(vaultPath, password, keyName, keyValue)
	if err != nil {
		t.Fatalf("KeySet failed: %v", err)
	}

	// Ensure the env var is not set before hydration
	os.Unsetenv(keyName)

	// Hydrate environment
	err = envenc.HydrateEnvFromFile(vaultPath, password)
	if err != nil {
		t.Fatalf("HydrateEnvFromFile failed: %v", err)
	}

	// Verify environment variable was set
	if os.Getenv(keyName) != keyValue {
		t.Errorf("Environment variable not set correctly: got %q, want %q", os.Getenv(keyName), keyValue)
	}

	// Cleanup
	os.Unsetenv(keyName)
}
