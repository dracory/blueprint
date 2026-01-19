package indexnow

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestIndexNowController_Handler_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetCmsStoreUsed(false) // use shared layout to avoid CMS dependencies
	app := testutils.Setup(testutils.WithCfg(cfg))

	controller := NewIndexNowController(app)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/indexnow", nil)
	html := controller.Handler(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Result().StatusCode)
	}

	if html == "" {
		t.Fatal("expected HTML to be non-empty")
	}

	if !strings.Contains(html, "IndexNow") {
		t.Fatalf("expected HTML to contain page title 'IndexNow'")
	}

	if !strings.Contains(html, "cd325dd195454606a8316fb303224f37") {
		t.Fatalf("expected HTML to contain the IndexNow key")
	}

	if !strings.Contains(html, "/cd325dd195454606a8316fb303224f37.txt") {
		t.Fatalf("expected HTML to contain link to key file")
	}
}
