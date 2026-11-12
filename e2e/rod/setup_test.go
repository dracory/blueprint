// Package rod contains end-to-end tests that drive the app with a real
// Chrome via go-rod (CDP). The e2e-browser sandbox server is started once per
// test run on http://127.0.0.1:34757.
package rod

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

const baseURL = "http://127.0.0.1:34757"

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

	code := m.Run()
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
