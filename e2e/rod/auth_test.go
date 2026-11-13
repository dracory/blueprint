package rod

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// Credentials of the users seeded by e2e/browser (seed.go).
const (
	seedAdminEmail = "e2e-admin@blueprint.local"
	seedUserEmail  = "e2e-user@blueprint.local"
	seedPassword   = "password123"
)

// // Launch a new browser with default options, and connect to it.
// 	browser := rod.New().MustConnect()

func newPage(t *testing.T) *rod.Page {
	t.Helper()

	// Leakless(false) avoids rod's leakless.exe helper, which Windows
	// Defender flags as a trojan (known false positive, go-rod#325).
	// Cleanup is handled by l.Kill in t.Cleanup instead.
	l := launcher.New().Headless(true).Leakless(false)
	if bin := chromeBin(); bin != "" {
		l = l.Bin(bin)
	}
	u := l.MustLaunch()
	t.Cleanup(l.Kill)

	browser := rod.New().ControlURL(u).MustConnect()
	t.Cleanup(func() { _ = browser.Close() })

	// Timeout returns a new page value — it must be assigned once here so
	// every subsequent call (Element, MustElement, ...) inherits it.
	return browser.MustPage().Timeout(15 * time.Second)
}

// chromeBin resolves a headless Chromium binary. Order: CHROME_BIN env,
// Playwright's downloaded headless shell, empty (rod auto-downloads a
// managed Chromium to ~/.cache/rod).
func chromeBin() string {
	if b := os.Getenv("CHROME_BIN"); b != "" {
		return b
	}
	patterns := []string{
		// Linux
		filepath.Join(os.Getenv("HOME"),
			".cache/ms-playwright/chromium_headless_shell-*/chrome-headless-shell-linux64/chrome-headless-shell"),
		// Windows
		filepath.Join(os.Getenv("LOCALAPPDATA"),
			"ms-playwright/chromium_headless_shell-*/chrome-headless-shell-win64/chrome-headless-shell.exe"),
	}
	for _, pattern := range patterns {
		if matches, _ := filepath.Glob(pattern); len(matches) > 0 {
			return matches[len(matches)-1]
		}
	}
	return ""
}

func TestLoginPageLoads(t *testing.T) {
	page := newPage(t)

	page.MustNavigate(baseURL + "/auth/login").MustWaitLoad()

	page.MustElement(`input[name="email"]`)
	page.MustElement(`input[name="password"]`)
}

func TestHomePageLoads(t *testing.T) {
	page := newPage(t)

	page.MustNavigate(baseURL + "/").MustWaitLoad()

	// The home controller renders "Welcome to <appName>"; the sandbox sets
	// APP_NAME=Blueprint.
	h1 := page.MustElement("h1").MustText()
	if !strings.Contains(strings.ToLower(h1), "blueprint") {
		t.Fatalf("expected heading to mention Blueprint, got %q", h1)
	}
}

func TestSeededAdminLandsOnAdmin(t *testing.T) {
	page := login(t, seedAdminEmail)
	waitForPath(t, page, "/admin")
}

func TestSeededUserLandsOnUser(t *testing.T) {
	page := login(t, seedUserEmail)
	waitForPath(t, page, "/user")
}

func TestSeededUserCannotAccessAdmin(t *testing.T) {
	page := login(t, seedUserEmail)
	waitForPath(t, page, "/user")

	// AdminMiddleware redirects non-admins to "/" with a flash error.
	// Wait via lifecycle events instead of MustWaitLoad: the /admin request
	// 303-redirects to a flash page, and evaluating JS mid-redirect makes
	// Chrome fail with "Object reference chain is too long".
	wait := page.MustWaitNavigation()
	page.MustNavigate(baseURL + "/admin")
	wait()

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		url := page.MustInfo().URL
		if !strings.Contains(url, "/admin") {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("regular user should not stay on /admin, still on %s", page.MustInfo().URL)
}

// login opens the login page, submits the seeded credentials, and returns
// the page for further assertions.
func login(t *testing.T, email string) *rod.Page {
	t.Helper()
	page := newPage(t)

	page.MustNavigate(baseURL + "/auth/login").MustWaitLoad()

	page.MustElement(`input[name="email"]`).MustInput(email)
	page.MustElement(`input[name="password"]`).MustInput(seedPassword)

	// The dracory/auth package submits via JS onclick, not a submit button:
	// <button class="ButtonLogin ...">. Try likely candidates in order.
	for _, sel := range []string{`button.ButtonLogin`, `button[type="submit"]`, `input[type="submit"]`} {
		if el, err := page.Element(sel); err == nil {
			el.MustClick()
			return page
		}
	}
	t.Fatal("no login submit button found")
	return nil
}

// waitForPath polls the page URL until it contains wantPath.
func waitForPath(t *testing.T, page *rod.Page, wantPath string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if strings.Contains(page.MustInfo().URL, wantPath) {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("expected redirect containing %q, still on %s", wantPath, page.MustInfo().URL)
}
