package auth

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"project/config"
	"project/internal/testutils"
	"testing"

	"github.com/dracory/base/test"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/vaultstore"
	_ "github.com/mattn/go-sqlite3"
)

func setupVaultStore(_ *testing.T) (vaultStore *vaultstore.Store, err error) {
	// Create a mock database connection
	mockDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// Initialize the vault store
	vaultStore, err = vaultstore.NewStore(vaultstore.NewStoreOptions{
		DB:                 mockDB,
		VaultTableName:     "snv_vault_vault",
		AutomigrateEnabled: true,
	})
	if err != nil {
		return nil, err
	}

	return vaultStore, nil
}

func TestLoginControllerHandler_UserStoreNotUsed(t *testing.T) {
	testutils.Setup()

	config.UserStoreUsed = false
	req, err := test.NewRequest(http.MethodGet, "/", test.NewRequestOptions{})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	(http.Handler(responses.HTMLHandler(NewLoginController().Handler))).ServeHTTP(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, recorder.Code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(recorder.Result())
	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	if flashMessage.Type != "error" {
		t.Fatal(`Response be of type 'error', but got: `, flashMessage.Type, flashMessage.Message)
	}

	if flashMessage.Message != "user store is required" {
		t.Fatal(`Response MUST contain 'user store is required', but got: `, flashMessage.Message)
	}
}

func TestLoginControllerHandler_VaultStoreNotUsed(t *testing.T) {
	testutils.Setup()

	config.UserStoreUsed = true
	config.VaultStoreUsed = true
	config.VaultStore = nil
	req, err := test.NewRequest(http.MethodGet, "/", test.NewRequestOptions{})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	(http.Handler(responses.HTMLHandler(NewLoginController().Handler))).ServeHTTP(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, recorder.Code)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(recorder.Result())
	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal(`Response MUST contain 'flash message'`)
	}

	if flashMessage.Type != "error" {
		t.Fatal(`Response be of type 'error', but got: `, flashMessage.Type, flashMessage.Message)
	}

	if flashMessage.Message != "vault store is required" {
		t.Fatal(`Response MUST contain 'vault store is required', but got: `, flashMessage.Message)
	}
}

func TestLoginControllerHandler_ValidRedirectWithBackURL(t *testing.T) {
	testutils.Setup()

	vaultStore, err := setupVaultStore(t)
	if err != nil {
		t.Fatal(err)
	}
	config.VaultStore = vaultStore
	config.VaultStoreUsed = true

	backUrl := "http://localhost/home"
	req, err := test.NewRequest(http.MethodGet, "/", test.NewRequestOptions{
		GetValues: url.Values{
			"back_url": {backUrl},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	(http.Handler(responses.HTMLHandler(NewLoginController().Handler))).ServeHTTP(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, recorder.Code)
	}
}

func TestLoginControllerHandler_InvalidBackURL(t *testing.T) {
	testutils.Setup()

	vaultStore, err := setupVaultStore(t)
	if err != nil {
		t.Fatal(err)
	}
	config.VaultStore = vaultStore
	config.VaultStoreUsed = true

	invalidBackUrl := "http://evil.com"
	req, err := test.NewRequest(http.MethodGet, "/", test.NewRequestOptions{
		GetValues: url.Values{
			"back_url": {invalidBackUrl},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	(http.Handler(responses.HTMLHandler(NewLoginController().Handler))).ServeHTTP(recorder, req)

	if recorder.Code != http.StatusSeeOther {
		t.Fatal(`Response MUST be 303`, recorder.Code)
	}
}
