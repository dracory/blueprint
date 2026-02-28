package blog

import (
	"bytes"
	"crypto/subtle"
	"fmt"
	"io"
	"net/http"
	"project/internal/links"
	"project/internal/registry"
	"strings"
	"time"

	"project/internal/controllers/website/blog/home"
	"project/internal/controllers/website/blog/post"

	blogstoreMcp "github.com/dracory/blogstore/mcp"
	"github.com/dracory/rtr"
)

type teeResponseWriter struct {
	w      http.ResponseWriter
	status int
	buf    bytes.Buffer
}

func (t *teeResponseWriter) Header() http.Header {
	return t.w.Header()
}

func (t *teeResponseWriter) WriteHeader(code int) {
	t.status = code
	t.w.WriteHeader(code)
}

func (t *teeResponseWriter) Write(p []byte) (int, error) {
	_, _ = t.buf.Write(p)
	return t.w.Write(p)
}

func (t *teeResponseWriter) Flush() {
	if f, ok := t.w.(http.Flusher); ok {
		f.Flush()
	}
}

func redactHeader(name string) bool {
	switch strings.ToLower(name) {
	case "authorization", "x-mcp-api-key", "cookie", "set-cookie":
		return true
	default:
		return false
	}
}

func logRequest(method string, path string, headers http.Header, body []byte) {
	const maxBody = 4096
	if len(body) > maxBody {
		body = body[:maxBody]
	}
	var sb strings.Builder
	sb.WriteString("[MCP DEBUG] request ")
	sb.WriteString(method)
	sb.WriteString(" ")
	sb.WriteString(path)
	sb.WriteString("\n")
	for k, vv := range headers {
		if redactHeader(k) {
			sb.WriteString(k)
			sb.WriteString(": [REDACTED]\n")
			continue
		}
		for _, v := range vv {
			sb.WriteString(k)
			sb.WriteString(": ")
			sb.WriteString(v)
			sb.WriteString("\n")
		}
	}
	if len(body) > 0 {
		sb.WriteString("body: ")
		sb.Write(body)
		sb.WriteString("\n")
	}
	fmt.Println(sb.String())
}

func logResponse(status int, body []byte) {
	const maxBody = 4096
	if len(body) > maxBody {
		body = body[:maxBody]
	}
	fmt.Printf("[MCP DEBUG] response status=%d body=%s", status, string(body))
}

func Routes(
	registry registry.RegistryInterface,
) []rtr.RouteInterface {
	mcpBlogHealthRoute := rtr.NewRoute().
		SetName("Website > Blog > MCP Endpoint > Health").
		SetPath(links.MCP_BLOG).
		SetMethod(http.MethodGet).
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			accept := r.Header.Get("Accept")
			if strings.Contains(accept, "text/event-stream") {
				rc := http.NewResponseController(w)

				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("event: ready\ndata: {}\n\n"))
				if err := rc.Flush(); err != nil {
					return
				}

				ticker := time.NewTicker(15 * time.Second)
				defer ticker.Stop()
				for {
					select {
					case <-r.Context().Done():
						return
					case <-ticker.C:
						_, _ = w.Write([]byte(": keepalive\n\n"))
						_ = rc.Flush()
					}
				}
			}

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("MCP is working"))
		})

	mcpBlogRoute := rtr.NewRoute().
		SetName("Website > Blog > MCP Endpoint").
		SetPath(links.MCP_BLOG).
		SetMethod(http.MethodPost).
		SetHandler(func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
			_ = r.Body.Close()
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			logRequest(r.Method, r.URL.Path, r.Header, bodyBytes)

			apiKey := r.Header.Get("X-MCP-API-Key")
			expectedKey := ""
			if registry != nil && registry.GetConfig() != nil {
				expectedKey = registry.GetConfig().GetCmsMcpApiKey()
			}
			if strings.TrimSpace(expectedKey) == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"error":"MCP API key not configured","message":"Set MCP_API_KEY environment variable"}`))
				return
			}
			if subtle.ConstantTimeCompare([]byte(apiKey), []byte(expectedKey)) != 1 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"Invalid API key","message":"Check X-MCP-API-Key header"}`))
				return
			}

			if registry == nil || registry.GetBlogStore() == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"Blog store not available","message":"Blog store is not initialized"}`))
				return
			}

			tw := &teeResponseWriter{w: w, status: http.StatusOK}
			m := blogstoreMcp.NewMCP(registry.GetBlogStore())
			m.Handler(tw, r)
			logResponse(tw.status, tw.buf.Bytes())
		})

	blogRoute := rtr.NewRoute().
		SetName("Guest > Blog").
		SetPath(links.BLOG).
		SetHTMLHandler(home.NewBlogController(registry).Handler)

	blogPostRegex01Route := rtr.NewRoute().
		SetName("Guest > Blog > Post with ID > Index").
		SetPath(links.BLOG_POST_WITH_REGEX).
		SetHTMLHandler(post.NewPostController(registry).Handler)

	blogPostRegex02Route := rtr.NewRoute().
		SetName("Guest > Blog > Post with ID && Title > Index").
		SetPath(links.BLOG_POST_WITH_REGEX2).
		SetHTMLHandler(post.NewPostController(registry).Handler)

	// blogPost01Route := rtr.NewRoute().
	// 	SetName("Guest > Blog > Post (ID)").
	// 	SetPath(links.BLOG_01).
	// 	SetHTMLHandler(post.NewPostController(registry).Handler)

	// blogPost02Route := rtr.NewRoute().
	// 	SetName("Guest > Blog > Post (ID && Title)").
	// 	SetPath(links.BLOG_02).
	// 	SetHTMLHandler(post.NewPostController(registry).Handler)

	blogPost01Route := rtr.NewRoute().
		SetName("Guest > Blog > Post (ID)").
		SetPath(links.BLOG_POST_01).
		SetHTMLHandler(post.NewPostController(registry).Handler)

	blogPost02Route := rtr.NewRoute().
		SetName("Guest > Blog > Post (ID && Title)").
		SetPath(links.BLOG_POST_02).
		SetHTMLHandler(post.NewPostController(registry).Handler)

	return []rtr.RouteInterface{
		mcpBlogHealthRoute,
		mcpBlogRoute,
		blogRoute,
		blogPostRegex01Route,
		blogPostRegex02Route,
		blogPost01Route,
		blogPost02Route,
	}
}
