package chromedp

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// navigate loads url in the tab and waits for the page body to be
// ready — the pages are server-rendered so the DOM is complete at
// that point.
func navigate(t *testing.T, ctx context.Context, url string) {
	t.Helper()
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
	)
	if err != nil {
		t.Fatalf("navigate %s failed: %v", url, err)
	}
}

// count returns the number of elements matching the CSS selector.
func count(t *testing.T, ctx context.Context, sel string) int {
	t.Helper()
	var n int
	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf(`document.querySelectorAll(%q).length`, sel), &n),
	)
	if err != nil {
		t.Fatalf("query %q failed: %v", sel, err)
	}
	return n
}

// mustExist asserts every selector matches at least one element.
func mustExist(t *testing.T, ctx context.Context, sels ...string) {
	t.Helper()
	for _, sel := range sels {
		if n := count(t, ctx, sel); n == 0 {
			t.Errorf("expected at least one element matching %q", sel)
		}
	}
}

// bodyText returns the rendered page text.
func bodyText(t *testing.T, ctx context.Context) string {
	t.Helper()
	var text string
	if err := chromedp.Run(ctx, chromedp.Evaluate(`document.body.innerText`, &text)); err != nil {
		t.Fatalf("reading body text failed: %v", err)
	}
	return text
}

// location returns the tab's current URL.
func location(t *testing.T, ctx context.Context) string {
	t.Helper()
	var loc string
	if err := chromedp.Run(ctx, chromedp.Location(&loc)); err != nil {
		t.Fatalf("reading location failed: %v", err)
	}
	return loc
}

// waitForPath polls the tab URL until it contains wantPath.
func waitForPath(t *testing.T, ctx context.Context, wantPath string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if strings.Contains(location(t, ctx), wantPath) {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("expected redirect containing %q, still on %s", wantPath, location(t, ctx))
}

// waitForPathGone polls the tab URL until it no longer contains path.
func waitForPathGone(t *testing.T, ctx context.Context, path string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if !strings.Contains(location(t, ctx), path) {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("expected to leave %q, still on %s", path, location(t, ctx))
}
