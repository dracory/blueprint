package chromedp

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// Credentials of the users seeded by e2e/browser (seed.go).
const (
	seedAdminEmail = "e2e-admin@blueprint.local"
	seedUserEmail  = "e2e-user@blueprint.local"
	seedPassword   = "password123"
)

func TestLoginPageLoads(t *testing.T) {
	ctx := newTab(t)
	navigate(t, ctx, baseURL+"/auth/login")
	mustExist(t, ctx, `input[name="email"]`, `input[name="password"]`)
}

func TestHomePageLoads(t *testing.T) {
	ctx := newTab(t)
	navigate(t, ctx, baseURL+"/")

	// The home controller renders "Welcome to <appName>"; the sandbox sets
	// APP_NAME=Blueprint.
	var h1 string
	err := chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelector('h1')?.innerText ?? ''`, &h1),
	)
	if err != nil {
		t.Fatalf("reading h1 failed: %v", err)
	}
	if !strings.Contains(strings.ToLower(h1), "blueprint") {
		t.Fatalf("expected heading to mention Blueprint, got %q", h1)
	}
}

func TestSeededAdminLandsOnAdmin(t *testing.T) {
	ctx := login(t, seedAdminEmail)
	waitForPath(t, ctx, "/admin")
}

func TestSeededUserLandsOnUser(t *testing.T) {
	ctx := login(t, seedUserEmail)
	waitForPath(t, ctx, "/user")
}

func TestSeededUserCannotAccessAdmin(t *testing.T) {
	ctx := login(t, seedUserEmail)
	waitForPath(t, ctx, "/user")

	// AdminMiddleware redirects non-admins away from /admin.
	navigate(t, ctx, baseURL+"/admin")
	waitForPathGone(t, ctx, "/admin")
}

// login opens the login page, submits the seeded credentials, and
// returns the tab context for further assertions. The page script
// waits 2s on success before redirecting, so callers must poll the
// location (waitForPath) rather than assert immediately.
func login(t *testing.T, email string) context.Context {
	t.Helper()
	ctx := newTab(t)
	err := chromedp.Run(ctx,
		chromedp.Navigate(baseURL+"/auth/login"),
		chromedp.WaitReady(`input[name="email"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="email"]`, email, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="password"]`, seedPassword, chromedp.ByQuery),
		// The dracory/auth package submits via JS onclick, not a submit
		// button: <button class="ButtonLogin ...">.
		chromedp.Click(`button.ButtonLogin`, chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("login as %s failed: %v", email, err)
	}
	return ctx
}
