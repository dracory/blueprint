// Package rod contains end-to-end tests that drive the app with a real
// Chrome via go-rod (CDP). The e2e-browser sandbox server is started once per
// test run on http://127.0.0.1:34757.
package rod

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const baseURL = "http://127.0.0.1:34757"

// sharedBrowser is launched once in TestMain; each test opens a fresh tab on
// it so tests don't pay full Chrome startup per test.
var sharedBrowser *rod.Browser

func TestMain(m *testing.M) {
	// Build the sandboxed app once; killing `go run` would only kill the
	// wrapper and leave the child server holding the DB.
	bin := "../../tmp/e2e-browser-bin"
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	build := exec.Command("go", "build", "-o", bin, "./e2e/browser")
	build.Dir = "../.."
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Println("failed to build e2e-browser:", err)
		os.Exit(1)
	}

	// Start it from the repository root so relative paths
	// (tmp/e2e-browser.db) resolve correctly.
	cmd := exec.Command(bin, "-reset")
	cmd.Dir = "../.."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("failed to start e2e-browser:", err)
		os.Exit(1)
	}

	if err := waitForServer(baseURL, 180*time.Second); err != nil {
		fmt.Println("e2e-browser did not become ready:", err)
		_ = cmd.Process.Kill()
		os.Exit(1)
	}

	// One shared browser process for the entire test run.
	// Leakless(false) avoids rod's leakless.exe helper, which Windows
	// Defender flags as a trojan (known false positive, go-rod#325).
	// Cleanup is handled by l.Kill below instead.
	l := launcher.New().Headless(true).Leakless(false)
	if bin := chromeBin(); bin != "" {
		l = l.Bin(bin)
	}
	u := l.MustLaunch()

	sharedBrowser = rod.New().ControlURL(u)
	if err := sharedBrowser.Connect(); err != nil {
		fmt.Println("failed to connect to browser:", err)
		l.Kill()
		_ = cmd.Process.Kill()
		os.Exit(1)
	}

	code := m.Run()

	_ = sharedBrowser.Close()
	l.Kill()
	_ = cmd.Process.Kill()
	os.Exit(code)
}

func waitForServer(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timeout after %s", timeout)
}

// chromeBin resolves a headless Chromium binary. Order: CHROME_BIN env,
// Playwright's downloaded headless shell, empty (rod auto-downloads a
// managed Chromium to ~/.cache/rod).
func chromeBin() string {
	if b := os.Getenv("CHROME_BIN"); b != "" {
		return b
	}
	patterns := []string{
		// Windows
		filepath.Join(os.Getenv("LOCALAPPDATA"),
			"ms-playwright/chromium_headless_shell-*/chrome-headless-shell-win64/chrome-headless-shell.exe"),
		// WSL: Windows-side Playwright shell via interop (the Linux shell
		// may be present but unusable due to missing libs like libnspr4)
		"/mnt/*/Users/*/AppData/Local/ms-playwright/chromium_headless_shell-*/chrome-headless-shell-win64/chrome-headless-shell.exe",
		// Linux
		filepath.Join(os.Getenv("HOME"),
			".cache/ms-playwright/chromium_headless_shell-*/chrome-headless-shell-linux64/chrome-headless-shell"),
	}
	for _, pattern := range patterns {
		if matches, _ := filepath.Glob(pattern); len(matches) > 0 {
			return matches[len(matches)-1]
		}
	}
	return ""
}
