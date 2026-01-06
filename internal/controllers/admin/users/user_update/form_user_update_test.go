package admin

import (
	"context"
	"net/url"
	"testing"

	"project/internal/ext"
	"project/internal/testutils"
	"project/internal/types"

	"github.com/dracory/userstore"
)

func setupAppAndUser(t *testing.T) (types.RegistryInterface, userstore.UserInterface) {
	t.Helper()

	app := testutils.Setup(
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
		testutils.WithGeoStore(true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	user.SetStatus(userstore.USER_STATUS_ACTIVE)
	user.SetMemo("Initial memo")
	user.SetCountry("GB")
	user.SetTimezone("Europe/London")

	if err := app.GetUserStore().UserUpdate(context.Background(), user); err != nil {
		t.Fatalf("UserUpdate returned error: %v", err)
	}

	tokenizeUserForTest(t, app, user, "John", "Doe", "john@example.com", "JD Consulting", "+44111222333")

	return app, user
}

func tokenizeUserForTest(t *testing.T, app types.RegistryInterface, user userstore.UserInterface, firstName, lastName, email, businessName, phone string) {
	t.Helper()

	ctx := context.Background()

	firstToken, lastToken, emailToken, phoneToken, businessToken, err := ext.UserTokenize(
		ctx,
		app.GetVaultStore(),
		app.GetConfig().GetVaultStoreKey(),
		user,
		firstName,
		lastName,
		email,
		phone,
		businessName,
	)
	if err != nil {
		t.Fatalf("UserTokenize returned error: %v", err)
	}

	user.SetFirstName(firstToken)
	user.SetLastName(lastToken)
	user.SetEmail(emailToken)
	user.SetPhone(phoneToken)
	user.SetBusinessName(businessToken)

	if err := app.GetUserStore().UserUpdate(ctx, user); err != nil {
		t.Fatalf("UserUpdate returned error: %v", err)
	}
}

func newMountedForm(t *testing.T, app types.RegistryInterface, user userstore.UserInterface, returnURL string) *formUserUpdate {
	t.Helper()

	component := NewFormUserUpdate(app)
	if component == nil {
		t.Fatalf("NewFormUserUpdate returned nil component")
	}

	form, ok := component.(*formUserUpdate)
	if !ok {
		t.Fatalf("component should be *formUserUpdate, got %T", component)
	}

	if err := form.Mount(context.Background(), map[string]string{
		"user_id":    user.ID(),
		"return_url": returnURL,
	}); err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	if form.FormError != "" {
		t.Fatalf("Mount produced form error: %s", form.FormError)
	}

	return form
}

func cloneValues(values url.Values) url.Values {
	clone := url.Values{}
	for key, list := range values {
		clone[key] = append([]string(nil), list...)
	}
	return clone
}

func TestFormUserUpdate_CountryChangeClearsTimezone(t *testing.T) {
	app, user := setupAppAndUser(t)
	form := newMountedForm(t, app, user, "/admin/users")

	if err := form.Handle(context.Background(), "country_change", url.Values{
		"user_country": {"US"},
	}); err != nil {
		t.Fatalf("country change returned error: %v", err)
	}

	if form.FormTimezone != "" {
		t.Fatalf("expected timezone to reset after country change, got %q", form.FormTimezone)
	}
}

func TestFormUserUpdate_Mount(t *testing.T) {
	app, user := setupAppAndUser(t)
	form := newMountedForm(t, app, user, "/admin/users")

	if form.FormFirstName != "John" {
		t.Fatalf("expected first name 'John', got %q", form.FormFirstName)
	}

	if form.FormLastName != "Doe" {
		t.Fatalf("expected last name 'Doe', got %q", form.FormLastName)
	}

	if form.FormEmail != "john@example.com" {
		t.Fatalf("expected email 'john@example.com', got %q", form.FormEmail)
	}

	if form.FormMemo != "Initial memo" {
		t.Fatalf("expected memo 'Initial memo', got %q", form.FormMemo)
	}

	if form.FormBusiness != "JD Consulting" {
		t.Fatalf("expected business name 'JD Consulting', got %q", form.FormBusiness)
	}

	if form.FormPhone != "+44111222333" {
		t.Fatalf("expected phone '+44111222333', got %q", form.FormPhone)
	}

	if form.FormCountry != "GB" {
		t.Fatalf("expected country 'GB', got %q", form.FormCountry)
	}

	if form.FormTimezone != "Europe/London" {
		t.Fatalf("expected timezone 'Europe/London', got %q", form.FormTimezone)
	}

	if form.ReturnURL != "/admin/users" {
		t.Fatalf("expected return URL '/admin/users', got %q", form.ReturnURL)
	}

	if len(form.StatusOptions) == 0 {
		t.Fatal("expected status options to be populated")
	}
}

func TestFormUserUpdate_HandleValidation(t *testing.T) {
	testCases := []struct {
		name          string
		mutate        func(url.Values)
		expectedError string
	}{
		{
			name: "missing status",
			mutate: func(v url.Values) {
				v.Set("user_status", "")
			},
			expectedError: "Status is required",
		},
		{
			name: "missing first name",
			mutate: func(v url.Values) {
				v.Set("user_first_name", "")
			},
			expectedError: "First name is required",
		},
		{
			name: "missing last name",
			mutate: func(v url.Values) {
				v.Set("user_last_name", "")
			},
			expectedError: "Last name is required",
		},
		{
			name: "missing email",
			mutate: func(v url.Values) {
				v.Set("user_email", "")
			},
			expectedError: "Email is required",
		},
		{
			name: "invalid email",
			mutate: func(v url.Values) {
				v.Set("user_email", "not-an-email")
			},
			expectedError: "Invalid email address",
		},
		{
			name: "missing country",
			mutate: func(v url.Values) {
				v.Set("user_country", "")
			},
			expectedError: "Country is required",
		},
		{
			name: "missing timezone",
			mutate: func(v url.Values) {
				v.Set("user_timezone", "")
			},
			expectedError: "Timezone is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app, user := setupAppAndUser(t)
			form := newMountedForm(t, app, user, "/admin/users")

			values := url.Values{
				"user_id":            {user.ID()},
				"user_status":        {userstore.USER_STATUS_ACTIVE},
				"user_first_name":    {"Jane"},
				"user_last_name":     {"Smith"},
				"user_email":         {"jane@example.com"},
				"user_business_name": {"Acme Co"},
				"user_phone":         {"+123"},
				"user_country":       {"GB"},
				"user_timezone":      {"Europe/London"},
				"user_memo":          {"Notes"},
			}

			payload := cloneValues(values)
			if tc.mutate != nil {
				tc.mutate(payload)
			}

			if err := form.Handle(context.Background(), "apply", payload); err != nil {
				t.Fatalf("Handle returned error: %v", err)
			}

			if form.FormError != tc.expectedError {
				t.Fatalf("expected error %q, got %q", tc.expectedError, form.FormError)
			}

			if form.FormSuccess != "" {
				t.Fatalf("expected no success message, got %q", form.FormSuccess)
			}
		})
	}
}

func TestFormUserUpdate_HandleApplySuccess(t *testing.T) {
	app, user := setupAppAndUser(t)
	form := newMountedForm(t, app, user, "/admin/users")

	payload := url.Values{
		"user_id":            {user.ID()},
		"user_status":        {userstore.USER_STATUS_INACTIVE},
		"user_first_name":    {"Alice"},
		"user_last_name":     {"Smith"},
		"user_email":         {"alice@example.com"},
		"user_business_name": {"Alice Consulting"},
		"user_phone":         {"+441234567"},
		"user_country":       {"GB"},
		"user_timezone":      {"Europe/London"},
		"user_memo":          {"Updated memo"},
	}

	if err := form.Handle(context.Background(), "apply", payload); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if form.FormError != "" {
		t.Fatalf("expected no form error, got %q", form.FormError)
	}

	if form.FormSuccess != "User saved successfully" {
		t.Fatalf("expected success message, got %q", form.FormSuccess)
	}

	if form.FormRedirectTo != "" {
		t.Fatalf("expected no redirect, got %q", form.FormRedirectTo)
	}

	if form.DisplayName != "Alice Smith" {
		t.Fatalf("expected display name 'Alice Smith', got %q", form.DisplayName)
	}

	updatedUser, err := app.GetUserStore().UserFindByID(context.Background(), user.ID())
	if err != nil {
		t.Fatalf("UserFindByID returned error: %v", err)
	}
	if updatedUser == nil {
		t.Fatal("UserFindByID returned nil user")
	}

	firstName, lastName, email, businessName, phone, err := ext.UserUntokenize(context.Background(), app, app.GetConfig().GetVaultStoreKey(), updatedUser)
	if err != nil {
		t.Fatalf("UserUntokenize returned error: %v", err)
	}

	if firstName != "Alice" {
		t.Fatalf("expected stored first name 'Alice', got %q", firstName)
	}

	if lastName != "Smith" {
		t.Fatalf("expected stored last name 'Smith', got %q", lastName)
	}

	if email != "alice@example.com" {
		t.Fatalf("expected stored email 'alice@example.com', got %q", email)
	}

	if businessName != "Alice Consulting" {
		t.Fatalf("expected business name 'Alice Consulting', got %q", businessName)
	}

	if phone != "+441234567" {
		t.Fatalf("expected phone '+441234567', got %q", phone)
	}

	if updatedUser.Status() != userstore.USER_STATUS_INACTIVE {
		t.Fatalf("expected status %q, got %q", userstore.USER_STATUS_INACTIVE, updatedUser.Status())
	}

	if updatedUser.Memo() != "Updated memo" {
		t.Fatalf("expected memo 'Updated memo', got %q", updatedUser.Memo())
	}

	if updatedUser.Country() != "GB" {
		t.Fatalf("expected country 'GB', got %q", updatedUser.Country())
	}

	if updatedUser.Timezone() != "Europe/London" {
		t.Fatalf("expected timezone 'Europe/London', got %q", updatedUser.Timezone())
	}
}

func TestFormUserUpdate_HandleSaveRedirect(t *testing.T) {
	app, user := setupAppAndUser(t)
	returnURL := "/admin/users/list"
	form := newMountedForm(t, app, user, returnURL)

	payload := url.Values{
		"user_id":            {user.ID()},
		"user_status":        {userstore.USER_STATUS_ACTIVE},
		"user_first_name":    {"Bob"},
		"user_last_name":     {"Taylor"},
		"user_email":         {"bob.taylor@example.com"},
		"user_business_name": {"BT Consulting"},
		"user_phone":         {"+358123"},
		"user_country":       {"FI"},
		"user_timezone":      {"Europe/Helsinki"},
		"user_memo":          {"Final memo"},
	}

	if err := form.Handle(context.Background(), "save", payload); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if form.FormError != "" {
		t.Fatalf("expected no form error, got %q", form.FormError)
	}

	if form.FormSuccess != "User saved successfully" {
		t.Fatalf("expected success message, got %q", form.FormSuccess)
	}

	if form.FormRedirectTo != returnURL {
		t.Fatalf("expected redirect to %q, got %q", returnURL, form.FormRedirectTo)
	}

	updatedUser, err := app.GetUserStore().UserFindByID(context.Background(), user.ID())
	if err != nil {
		t.Fatalf("UserFindByID returned error: %v", err)
	}
	if updatedUser == nil {
		t.Fatal("UserFindByID returned nil user")
	}

	firstName, lastName, email, businessName, phone, err := ext.UserUntokenize(context.Background(), app, app.GetConfig().GetVaultStoreKey(), updatedUser)
	if err != nil {
		t.Fatalf("UserUntokenize returned error: %v", err)
	}

	if firstName != "Bob" {
		t.Fatalf("expected stored first name 'Bob', got %q", firstName)
	}

	if lastName != "Taylor" {
		t.Fatalf("expected stored last name 'Taylor', got %q", lastName)
	}

	if email != "bob.taylor@example.com" {
		t.Fatalf("expected stored email 'bob.taylor@example.com', got %q", email)
	}

	if businessName != "BT Consulting" {
		t.Fatalf("expected business name 'BT Consulting', got %q", businessName)
	}

	if phone != "+358123" {
		t.Fatalf("expected phone '+358123', got %q", phone)
	}

	if updatedUser.Memo() != "Final memo" {
		t.Fatalf("expected memo 'Final memo', got %q", updatedUser.Memo())
	}

	if updatedUser.Country() != "FI" {
		t.Fatalf("expected country 'FI', got %q", updatedUser.Country())
	}

	if updatedUser.Timezone() != "Europe/Helsinki" {
		t.Fatalf("expected timezone 'Europe/Helsinki', got %q", updatedUser.Timezone())
	}
}
