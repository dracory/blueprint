package account

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"

	"project/internal/ext"
	"project/internal/links"
	"project/internal/testutils"
	"project/internal/types"

	"github.com/dracory/userstore"
)

func userTokenize(app types.AppInterface, user userstore.UserInterface) error {
	ctx := context.Background()
	vaultStore := app.GetVaultStore()
	vaultKey := app.GetConfig().GetVaultStoreKey()

	ensureToken := func(existingToken string, value string, field string) (string, error) {
		if existingToken == "" {
			token, err := vaultStore.TokenCreate(ctx, value, vaultKey, 20)
			if err != nil {
				return "", errors.New("TokenCreate (" + field + ") returned error: " + err.Error())
			}
			return token, nil
		}

		if err := vaultStore.TokenUpdate(ctx, existingToken, value, vaultKey); err != nil {
			return "", errors.New("TokenUpdate (" + field + ") returned error: " + err.Error())
		}

		return existingToken, nil
	}

	emailToken, err := ensureToken(user.Email(), "john@example.com", "email")
	if err != nil {
		return err
	}

	firstNameToken, err := ensureToken(user.FirstName(), "John", "first name")
	if err != nil {
		return err
	}

	lastNameToken, err := ensureToken(user.LastName(), "Doe", "last name")
	if err != nil {
		return err
	}

	businessNameToken, err := ensureToken(user.BusinessName(), "JD Consulting", "business name")
	if err != nil {
		return err
	}

	phoneToken, err := ensureToken(user.Phone(), "+44111222333", "phone")
	if err != nil {
		return err
	}

	user.SetEmail(emailToken)
	user.SetFirstName(firstNameToken)
	user.SetLastName(lastNameToken)
	user.SetBusinessName(businessNameToken)
	user.SetPhone(phoneToken)
	user.SetCountry("GB")
	user.SetTimezone("Europe/London")

	if err := app.GetUserStore().UserUpdate(ctx, user); err != nil {
		return errors.New("UserUpdate returned error: " + err.Error())
	}

	return nil
}

func TestFormProfileUpdate_Mount(t *testing.T) {
	app := testutils.Setup(
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	if err := userTokenize(app, user); err != nil {
		t.Fatalf("TokenUpdate (email) returned error: %v", err)
	}

	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id": user.ID(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	if form.FormError != "" {
		t.Fatalf("Expected no error, got: %s", form.FormError)
	}

	if form.FormFirstName != "John" {
		t.Fatalf("Expected first name 'John', got: %s", form.FormFirstName)
	}

	if form.FormCountry == "" {
		t.Fatal("Expected country to be populated")
	}

	if len(form.Countries) == 0 {
		t.Fatal("Expected Countries to be loaded")
	}

	if len(form.Timezones) == 0 {
		t.Fatal("Expected Timezones to be loaded")
	}
}

func TestFormProfileUpdate_Handle_RequiresFirstName(t *testing.T) {
	app := testutils.Setup(
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id": user.ID(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	formData := url.Values{
		"user_id":    {user.ID()},
		"email":      {user.Email()},
		"first_name": {""},
		"last_name":  {"LastName"},
		"country":    {"Country"},
		"timezone":   {"Timezone"},
	}

	err = form.Handle(context.Background(), "apply", formData)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if form.FormError != "First name is required field" {
		t.Fatalf("Expected error %q, got: %q", "First name is required field", form.FormError)
	}

	if form.FormSuccess != "" {
		t.Fatalf("Expected no success message, got: %s", form.FormSuccess)
	}
}

func TestFormProfileUpdate_Handle_RequiresCountry(t *testing.T) {
	app := testutils.Setup(
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	if err := userTokenize(app, user); err != nil {
		t.Fatalf("TokenizeUser returned error: %v", err)
	}

	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id": user.ID(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	formData := url.Values{
		"user_id":    {user.ID()},
		"email":      {"user@example.com"},
		"first_name": {"FirstName"},
		"last_name":  {"LastName"},
		"country":    {""},
		"timezone":   {"Timezone"},
	}

	err = form.Handle(context.Background(), "apply", formData)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if form.FormError != "Country is required field" {
		t.Fatalf("Expected error %q, got: %q", "Country is required field", form.FormError)
	}

	if form.FormSuccess != "" {
		t.Fatalf("Expected no success message, got: %s", form.FormSuccess)
	}
}

func TestFormProfileUpdate_Handle_RequiresTimezone(t *testing.T) {
	app := testutils.Setup(
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	if err := userTokenize(app, user); err != nil {
		t.Fatalf("TokenizeUser returned error: %v", err)
	}

	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id": user.ID(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	formData := url.Values{
		"user_id":    {user.ID()},
		"email":      {"user@example.com"},
		"first_name": {"FirstName"},
		"last_name":  {"LastName"},
		"country":    {"Country"},
		"timezone":   {""},
	}

	err = form.Handle(context.Background(), "apply", formData)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if form.FormError != "Timezone is required field" {
		t.Fatalf("Expected error %q, got: %q", "Timezone is required field", form.FormError)
	}

	if form.FormSuccess != "" {
		t.Fatalf("Expected no success message, got: %s", form.FormSuccess)
	}
}

func TestFormProfileUpdate_Handle_RequiresLastName(t *testing.T) {
	app := testutils.Setup(
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	if err := userTokenize(app, user); err != nil {
		t.Fatalf("TokenizeUser returned error: %v", err)
	}

	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id": user.ID(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	formData := url.Values{
		"user_id":    {user.ID()},
		"email":      {user.Email()},
		"first_name": {"FirstName"},
		"last_name":  {""},
		"country":    {"Country"},
		"timezone":   {"Timezone"},
	}

	err = form.Handle(context.Background(), "apply", formData)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if form.FormError != "Last name is required field" {
		t.Fatalf("Expected error %q, got: %q", "Last name is required field", form.FormError)
	}

	if form.FormSuccess != "" {
		t.Fatalf("Expected no success message, got: %s", form.FormSuccess)
	}
}

func TestFormProfileUpdate_Handle_RequiresEmail(t *testing.T) {
	app := testutils.Setup(
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	if err := userTokenize(app, user); err != nil {
		t.Fatalf("TokenizeUser returned error: %v", err)
	}

	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id": user.ID(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	formData := url.Values{
		"user_id":    {user.ID()},
		"email":      {""},
		"first_name": {"FirstName"},
		"last_name":  {"LastName"},
		"country":    {"Country"},
		"timezone":   {"Timezone"},
	}

	err = form.Handle(context.Background(), "apply", formData)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if form.FormError != "Email is required field" {
		t.Fatalf("Expected error %q, got: %q", "Email is required field", form.FormError)
	}

	if form.FormSuccess != "" {
		t.Fatalf("Expected no success message, got: %s", form.FormSuccess)
	}
}

func TestFormProfileUpdate_Handle_Validation(t *testing.T) {
	app := testutils.Setup(
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	if err := userTokenize(app, user); err != nil {
		t.Fatalf("TokenizeUser returned error: %v", err)
	}

	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id": user.ID(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	testCases := []struct {
		name          string
		formData      url.Values
		expectedError string
	}{
		{
			name: "missing user id",
			formData: url.Values{
				"email":      {"jane@example.com"},
				"first_name": {"Jane"},
				"last_name":  {"Smith"},
				"country":    {"GB"},
				"timezone":   {"Europe/London"},
			},
			expectedError: "User ID is required",
		},
		{
			name: "missing first name",
			formData: url.Values{
				"user_id":    {user.ID()},
				"email":      {"jane@example.com"},
				"first_name": {""},
				"last_name":  {"Smith"},
				"country":    {"GB"},
				"timezone":   {"Europe/London"},
			},
			expectedError: "First name is required field",
		},
		{
			name: "missing last name",
			formData: url.Values{
				"user_id":    {user.ID()},
				"email":      {"jane@example.com"},
				"first_name": {"Jane"},
				"last_name":  {""},
				"country":    {"GB"},
				"timezone":   {"Europe/London"},
			},
			expectedError: "Last name is required field",
		},
		{
			name: "missing email",
			formData: url.Values{
				"user_id":    {user.ID()},
				"email":      {""},
				"first_name": {"Jane"},
				"last_name":  {"Smith"},
				"country":    {"GB"},
				"timezone":   {"Europe/London"},
			},
			expectedError: "Email is required field",
		},
		{
			name: "missing country",
			formData: url.Values{
				"user_id":    {user.ID()},
				"email":      {"jane@example.com"},
				"first_name": {"Jane"},
				"last_name":  {"Smith"},
				"country":    {""},
				"timezone":   {"Europe/London"},
			},
			expectedError: "Country is required field",
		},
		{
			name: "missing timezone",
			formData: url.Values{
				"user_id":    {user.ID()},
				"email":      {"jane@example.com"},
				"first_name": {"Jane"},
				"last_name":  {"Smith"},
				"country":    {"GB"},
				"timezone":   {""},
			},
			expectedError: "Timezone is required field",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := form.Handle(context.Background(), "apply", tc.formData)
			if err != nil {
				t.Fatalf("Handle returned error: %v", err)
			}

			if form.FormError != tc.expectedError {
				t.Fatalf("Expected error %q, got: %q", tc.expectedError, form.FormError)
			}

			if form.FormSuccess != "" {
				t.Fatalf("Expected no success message, got: %s", form.FormSuccess)
			}
		})
	}
}

func TestFormProfileUpdate_Handle_Apply(t *testing.T) {
	app := testutils.Setup(
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	if err := userTokenize(app, user); err != nil {
		t.Fatalf("TokenizeUser returned error: %v", err)
	}

	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id": user.ID(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	formData := url.Values{
		"user_id":       {user.ID()},
		"email":         {"john@example.com"},
		"first_name":    {"Johnny"},
		"last_name":     {"D"},
		"business_name": {"JD Consulting"},
		"phone":         {"+44111222333"},
		"country":       {"US"},
		"timezone":      {"America/New_York"},
	}

	err = form.Handle(context.Background(), "apply", formData)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if form.FormError != "" {
		t.Fatalf("Expected no error, got: %s", form.FormError)
	}

	if form.FormSuccess != "Profile updated successfully" {
		t.Fatalf("Expected success message, got: %s", form.FormSuccess)
	}

	if form.FormRedirectTo != "" {
		t.Fatalf("Expected no redirect for apply action, got: %s", form.FormRedirectTo)
	}

	// Verify the user was updated in the store
	updatedUser, err := app.GetUserStore().UserFindByID(context.Background(), user.ID())
	if err != nil {
		t.Fatalf("UserFindByID returned error: %v", err)
	}

	if updatedUser == nil {
		t.Fatalf("Expected user to be found")
	}

	firstName, _, _, _, _, err := ext.UserUntokenize(context.Background(), app, app.GetConfig().GetVaultStoreKey(), updatedUser)
	if err != nil {
		t.Fatalf("UserUntokenized returned error: %v", err)
	}

	if firstName != "Johnny" {
		t.Fatalf("Expected first name to be 'Johnny', got: %s", firstName)
	}

	if updatedUser.Country() != "US" {
		t.Fatalf("Expected country to be 'US', got: %s", updatedUser.Country())
	}

	if updatedUser.Timezone() != "America/New_York" {
		t.Fatalf("Expected timezone to be 'America/New_York', got: %s", updatedUser.Timezone())
	}
}

func TestFormProfileUpdate_Handle_Save_NoVault(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	user.SetEmail("user@example.com")
	if err := app.GetUserStore().UserUpdate(context.Background(), user); err != nil {
		t.Fatalf("UserUpdate returned error: %v", err)
	}

	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id":    user.ID(),
		"return_url": links.User().Profile(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	formData := url.Values{
		"user_id":       {user.ID()},
		"email":         {user.Email()},
		"first_name":    {"FirstName"},
		"last_name":     {"LastName"},
		"business_name": {"Biz"},
		"phone":         {"123"},
		"country":       {"Country"},
		"timezone":      {"Timezone"},
	}

	err = form.Handle(context.Background(), "save", formData)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if form.FormError != "" {
		t.Fatalf("Expected no error, got: %s", form.FormError)
	}

	if form.FormSuccess != "Profile updated successfully" {
		t.Fatalf("Expected success message, got: %s", form.FormSuccess)
	}

	if form.FormRedirectTo != links.User().Home() {
		t.Fatalf("Expected redirect to %s, got: %s", links.User().Home(), form.FormRedirectTo)
	}

	updatedUser, err := app.GetUserStore().UserFindByID(context.Background(), user.ID())
	if err != nil {
		t.Fatalf("UserFindByID returned error: %v", err)
	}

	if updatedUser == nil {
		t.Fatalf("Expected user to be found")
	}

	if updatedUser.Country() != "Country" {
		t.Fatalf("Expected country to be 'Country', got: %s", updatedUser.Country())
	}

	if updatedUser.Timezone() != "Timezone" {
		t.Fatalf("Expected timezone to be 'Timezone', got: %s", updatedUser.Timezone())
	}
}

func TestFormProfileUpdate_Handle_Save_WithVault(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	if err := userTokenize(app, user); err != nil {
		t.Fatalf("TokenizeUser returned error: %v", err)
	}

	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id":    user.ID(),
		"return_url": links.User().Profile(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	formData := url.Values{
		"user_id":       {user.ID()},
		"email":         {"user@example.com"},
		"first_name":    {"FirstName"},
		"last_name":     {"LastName"},
		"business_name": {"Biz"},
		"phone":         {"123"},
		"country":       {"Country"},
		"timezone":      {"Timezone"},
	}

	err = form.Handle(context.Background(), "save", formData)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if form.FormError != "" {
		t.Fatalf("Expected no error, got: %s", form.FormError)
	}

	if form.FormSuccess != "Profile updated successfully" {
		t.Fatalf("Expected success message, got: %s", form.FormSuccess)
	}

	if form.FormRedirectTo != links.User().Home() {
		t.Fatalf("Expected redirect to %s, got: %s", links.User().Home(), form.FormRedirectTo)
	}

	updatedUser, err := app.GetUserStore().UserFindByID(context.Background(), user.ID())
	if err != nil {
		t.Fatalf("UserFindByID returned error: %v", err)
	}

	if updatedUser == nil {
		t.Fatalf("Expected user to be found")
	}

	if updatedUser.Country() != "Country" {
		t.Fatalf("Expected country to be 'Country', got: %s", updatedUser.Country())
	}

	if updatedUser.Timezone() != "Timezone" {
		t.Fatalf("Expected timezone to be 'Timezone', got: %s", updatedUser.Timezone())
	}

	firstName, _, _, _, _, err := ext.UserUntokenize(context.Background(), app, app.GetConfig().GetVaultStoreKey(), updatedUser)
	if err != nil {
		t.Fatalf("UserUntokenized returned error: %v", err)
	}

	if firstName != "FirstName" {
		t.Fatalf("Expected first name to be 'FirstName', got: %s", firstName)
	}
}

func TestFormProfileUpdate_Handle_Save(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	if err := userTokenize(app, user); err != nil {
		t.Fatalf("TokenizeUser returned error: %v", err)
	}

	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id": user.ID(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	form.ReturnURL = "/user/home"

	formData := url.Values{
		"user_id":       {user.ID()},
		"email":         {"info@example.com"},
		"first_name":    {"Updated"},
		"last_name":     {"User"},
		"business_name": {"Updated Biz"},
		"phone":         {"+44111222331"},
		"country":       {"US"},
		"timezone":      {"America/New_York"},
	}

	err = form.Handle(context.Background(), "save", formData)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if form.FormError != "" {
		t.Fatalf("Expected no error, got: %s", form.FormError)
	}

	if form.FormRedirectTo == "" {
		t.Fatal("Expected redirect for save action")
	}

	if form.FormSuccess != "Profile updated successfully" {
		t.Fatalf("Expected success message 'Profile updated successfully', got: %s", form.FormSuccess)
	}

	// Verify the user was updated in the store
	updatedUser, err := app.GetUserStore().UserFindByID(context.Background(), user.ID())
	if err != nil {
		t.Fatalf("UserFindByID returned error: %v", err)
	}

	if updatedUser == nil {
		t.Fatalf("Expected user to be found")
	}

	firstName, _, _, _, _, err := ext.UserUntokenize(context.Background(), app, app.GetConfig().GetVaultStoreKey(), updatedUser)
	if err != nil {
		t.Fatalf("UserUntokenized returned error: %v", err)
	}

	if firstName != "Updated" {
		t.Fatalf("Expected first name to be 'Updated', got: %s", firstName)
	}

	if updatedUser.Country() != "US" {
		t.Fatalf("Expected country to be 'US', got: %s", updatedUser.Country())
	}

	if updatedUser.Timezone() != "America/New_York" {
		t.Fatalf("Expected timezone to be 'America/New_York', got: %s", updatedUser.Timezone())
	}
}

func TestFormProfileUpdate_Render(t *testing.T) {
	app := testutils.Setup(
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	if err := userTokenize(app, user); err != nil {
		t.Fatalf("TokenizeUser returned error: %v", err)
	}
	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id": user.ID(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	html := form.Render(context.Background())
	if html == nil {
		t.Fatal("Expected HTML to be rendered, got nil")
	}

	htmlStr := html.ToHTML()
	if htmlStr == "" {
		t.Fatal("Expected non-empty HTML string")
	}

	expecteds := []string{
		"Your Details",
		"First name",
		"Last name",
		"Email",
		"Country",
		"Timezone",
		"Save",
	}

	for _, expected := range expecteds {
		if !strings.Contains(htmlStr, expected) {
			t.Fatalf("Response MUST contain %q", expected)
		}
	}
}

func TestFormProfileUpdate_GetAlias(t *testing.T) {
	app := testutils.Setup()
	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	alias := form.GetAlias()

	if alias != "profile_update_form" {
		t.Fatalf("Expected alias 'profile_update_form', got: %s", alias)
	}
}

func TestFormProfileUpdate_Handle_CountryChange(t *testing.T) {
	app := testutils.Setup(
		testutils.WithGeoStore(true),
		testutils.WithUserStore(true, true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	component := NewFormProfileUpdate(app)
	form := component.(*formProfileUpdate)

	err = form.Mount(context.Background(), map[string]string{
		"user_id": user.ID(),
	})
	if err != nil {
		t.Fatalf("Mount returned error: %v", err)
	}

	// Simulate a country change
	formData := url.Values{
		"country": {"US"},
	}

	err = form.Handle(context.Background(), "country_change", formData)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if form.FormCountry != "US" {
		t.Fatalf("Expected country to be 'US', got: %s", form.FormCountry)
	}

	// Check if timezones were refreshed for the new country
	if len(form.Timezones) == 0 {
		t.Fatal("Expected Timezones to be refreshed and not empty for US")
	}

	// A more specific check could be to see if a known US timezone is present
	found := false
	for _, tz := range form.Timezones {
		if strings.HasPrefix(tz.Timezone(), "America/") {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Expected to find an 'America/' timezone for country US")
	}
}
