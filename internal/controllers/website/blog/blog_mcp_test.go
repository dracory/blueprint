package blog

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"project/internal/links"
	"project/internal/testutils"

	"github.com/dracory/rtr"
)

func TestBlogMcpEndpoint_RequiresApiKey(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetBlogStoreUsed(true)
	cfg.SetCmsMcpApiKey("test-mcp-key")

	registry := testutils.Setup(testutils.WithCfg(cfg))

	r := rtr.NewRouter()
	r.AddRoutes(Routes(registry))

	// Missing key
	reqMissing := httptest.NewRequest(http.MethodPost, links.MCP_BLOG, bytes.NewBuffer([]byte(`{"jsonrpc":"2.0","id":"1","method":"list_tools","params":{}}`)))
	resMissing := httptest.NewRecorder()
	r.ServeHTTP(resMissing, reqMissing)
	if resMissing.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, resMissing.Code)
	}

	// Wrong key
	reqWrong := httptest.NewRequest(http.MethodPost, links.MCP_BLOG, bytes.NewBuffer([]byte(`{"jsonrpc":"2.0","id":"1","method":"list_tools","params":{}}`)))
	reqWrong.Header.Set("X-MCP-API-Key", "wrong")
	resWrong := httptest.NewRecorder()
	r.ServeHTTP(resWrong, reqWrong)
	if resWrong.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, resWrong.Code)
	}

	// Correct key
	reqOk := httptest.NewRequest(http.MethodPost, links.MCP_BLOG, bytes.NewBuffer([]byte(`{"jsonrpc":"2.0","id":"1","method":"list_tools","params":{}}`)))
	reqOk.Header.Set("X-MCP-API-Key", "test-mcp-key")
	resOk := httptest.NewRecorder()
	r.ServeHTTP(resOk, reqOk)
	if resOk.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, resOk.Code)
	}
}
